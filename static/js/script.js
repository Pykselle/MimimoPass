function ViewOrHide(doc) {
    if (doc.className.indexOf("w3-show") == -1) {
        if (doc.className.indexOf("w3-hide") == -1) {
            doc.className += " w3-hide";
        } else {
            doc.className = doc.className.replace("w3-hide", "w3-show");
        }
    } else {
        doc.className = doc.className.replace("w3-show", "w3-hide");
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