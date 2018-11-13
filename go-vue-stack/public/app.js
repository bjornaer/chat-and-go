new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        username: null, // Our username
        joined: false // True if email and username have been filled in
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.parseMessage(msg, true);
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        content: $('<p>').html(this.newMsg).text() // Strip out html
                        //timestamp: new Date().toLocaleString()
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },

        join: function () {
            var self = this;
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            var email = $('<p>').html(this.email).text();
            var username = $('<p>').html(this.username).text();
            var opts = {"email": email, "username": username};
            fetch('/login', {
                method: 'post',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(opts)
            }).then(function(response) {
                if(response.status === 409 || response.status === 502) {
                    throw "Wrong username or email account!"
                }
                return response.json();
            }).then(function(data) {
                self.email = email; // $('<p>').html(this.email).text();
                self.username = username; // $('<p>').html(this.username).text();
                self.joined = true;
                self.historical();
                self.smoothScrollToBottom("chat-messages")
            })
            .catch(e => {
                Materialize.toast(e, 2000)
                console.log(e);
                return
            });
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        },

        historical: function() {
            // fetch /history endpoint and forEach do the parsing
            // TODO receive paginated result
            var oldestId;
            fetch(`/history?oldest=${oldestId}`)
            .then( response => {
                if(response.status !== 200) {
                    console.log('Whoops! Not the expected status! Status:' + response.status);
                    return
                }
                response.json()
                .then( data => {
                    var messages = data.messages.reverse()
                    messages.forEach((msg) => {
                        this.parseMessage(msg, false);
                    })
                })
            })
            .catch(err => {
                console.log('Error retrieving history: -S', err)
            });
        },

        unreadHist: function() {
            // fetch /news endpoint and forEach do the parsing
            var since;
            fetch(`/newMessages?since=${since}`)
            .then( response => {
                if(response.status !== 200) {
                    console.log('Whoops! Not the expected status! Status:' + response.status);
                    return
                }
                response.json()
                .then( data => {
                    var messages = data.messages.reverse()
                    messages.forEach((msg) => {
                        this.parseMessage(msg, false);
                    })
                })
            })
            .catch(err => {
                console.log('Error retrieving history: -S', err)
            });
        },

        parseMessage: function(msg, scroll) {
            this.chatContent += '<div class="chip">'
                    + '<img src="' + this.gravatarURL(msg.email) + '">' // Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.content) + '<br/>'; // Parse emojis
            if(scroll) {
                this.smoothScrollToBottom("chat-messages");
            }
        },

        smoothScrollToBottom: function(id) {
            var div = document.getElementById(id);
            $('#' + id).animate({
               scrollTop: div.scrollHeight - div.clientHeight + 500
            }, 500);
         }
    }
});
