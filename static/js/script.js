function ViewOrHide(doc) {
    if (doc.className.indexOf("w3-show") == -1) {
        doc.className += " w3-show";
    } else {
        doc.className = doc.className.replace(" w3-show", "");
    }
}

function copy() {
    let initial_state = document.getElementById('pass').type;
    if (initial_state != "text") {
        document.getElementById('pass').type = 'text';
        document.getElementById('pass').select();
        document.execCommand('copy');
        document.getElementById('pass').type = initial_state;
    } else {
        document.getElementById('pass').select();
        document.execCommand('copy');
    }

}