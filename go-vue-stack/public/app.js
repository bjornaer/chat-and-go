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
            self.chatContent += '<div class="chip">'
                    + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.content) + '<br/>'; // Parse emojis

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        content: $('<p>').html(this.newMsg).text(), // Strip out html
                        timestamp: new Date().toLocaleString()
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },

        join: function () {
            var _this = this;
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
                if(response.status === 409) {
                    Materialize.toast('Wrong username for email account', 2000)
                    return
                }
                return response.json();
            }).then(function(data) {
                _this.email = email; // $('<p>').html(this.email).text();
                _this.username = username; // $('<p>').html(this.username).text();
                _this.joined = true;
            })
            .catch(e => {
                console.log(e);
                return
            });
           this.historical();
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        },

        historical: function() {
            // fetch /history endpoint and forEach do the parsing
            fetch('/history')
            .then( response => {
                if(response.status !== 200) {
                    console.log('Whoops! Not the expected status! Status:' + response.status);
                    return
                }
                response.json()
                .then( data => {
                    var messages = data.messages.reverse()
                    messages.forEach((msg) => {
                        this.parseMessage(msg);
                    })
                })
            })
            .catch(err => {
                console.log('Error retrieving history: -S', err)
            });
        },

        parseMessage: function(msg) {
            this.chatContent += '<div class="chip">'
                    + '<img src="' + this.gravatarURL(msg.email) + '">' // Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.content) + '<br/>'; // Parse emojis

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        }
    }
});
