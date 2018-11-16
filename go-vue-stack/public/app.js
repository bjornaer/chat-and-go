new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        user: {
            email: null,
            username: null,
            id: null
        },
        oldestMessage: -1,
        joined: false // True if email and username have been filled in
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.parseMessage(msg, true, "bottom");
        });
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        userId: this.user.id,
                        content: $('<p>').html(this.newMsg).text() // Strip out html
                    }
                ));
                this.newMsg = ''; // Reset newMsg
            }
        },

        join: function () {
            var self = this;
            if (!this.user.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.user.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            var email = $('<p>').html(this.user.email).text();
            var username = $('<p>').html(this.user.username).text();
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
                self.user.id = data.id;
                self.user.username = data.username;
                self.user.email = data.email;
                self.joined = true;
                self.historical();
                self.smoothScrollToBottom("chat-messages")
            })
            .catch(e => {
                Materialize.toast("Review email and username",2000)
                console.log(e);
                return
            });
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        },

        historical: function() {
            if(this.oldestMessage === 1){
                return
            }
            fetch(`/history?oldest=${this.oldestMessage}&quantity=${50}`)
            .then( response => {
                if(response.status !== 200) {
                    Materialize.toast("Failed retrieving chat history!", 2000)
                    console.log('Whoops! Not the expected status! Status:' + response.status);
                    return
                }
                response.json()
                .then( data => {
                    var messages = data.messages;
                    this.oldestMessage = messages[messages.length - 1].id;
                    messages.forEach((msg) => {
                        this.parseMessage(msg, false, "top");
                    })
                })
            })
            .catch(err => {
                console.log('Error retrieving history: -S', err)
            });
        },

        unreadHistory: function() {
            fetch(`/newMessages?id=${this.user.id}`)
            .then( response => {
                if(response.status !== 200) {
                    console.log('Whoops! Not the expected status! Status:' + response.status);
                    return
                }
                response.json()
                .then( data => {
                    var messages = data.messages.reverse()
                    messages.forEach((msg) => {
                        this.parseMessage(msg, false, "bottom");
                    })
                })
            })
            .catch(err => {
                console.log('Error retrieving history: -S', err)
            });
        },

        parseMessage: function(msg, scroll, appendTo) {
            var messageElement = '<div class="chip">'
            + '<img src="' + this.gravatarURL(msg.user.email) + '">' // Avatar
            + msg.user.username
        + '</div>'
        + emojione.toImage(msg.content) + '<br/>'; // Parse emojis

            if(appendTo === "top") {
                this.chatContent = messageElement + this.chatContent;
            } else if(appendTo === "bottom") {
                this.chatContent += messageElement;
            } else {
                console.error("parseMessage needs appendTo parameter to be in [top, bottom]")
            }
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
