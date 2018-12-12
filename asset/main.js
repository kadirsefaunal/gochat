$(document).ready(function () {
    var socket = io.connect('http://localhost:5000');

    socket.on('message', function (message) {
        if ("username" in message) {
            $('#messages').append("<b>" + message.username + " " + message.message + "</b>")
        } else {
            $('#messages').append("<p><b>" + message.from + ": </b>" + message.message + "</p>");
        }
    });
    
    $('#connect').click(function () {
        var username = $("#username").val();
        socket.emit("set-user", username);    
    });

    $('#send').click(function () {
        message = {
            from: $("#username").val(),
            to: $('#to').val(),
            message: $('#msg').val()
        };

        socket.emit('message', JSON.stringify(message));
        $('#messages').append("<p><b>" + message.from + ": </b>" + message.message + "</p>");
    });

    socket.on('user-list', function (users) {
        $.each(users, function (i, val) {
            if ($('#username').val() !== val) {
                $('#user-list').append('<li>' + val + '</li>');    
            }
        });
    });
});