function ViewOrHide(doc) {
    if (doc.style.display === 'none') {
        doc.style.display = 'table';
    } else {
        doc.style.display = 'none';
    }
}

function ViewOrHidePass() {
    let docPass = document.getElementById('pass')
    let docIcon = document.getElementById('viewIcon')
    if (docPass.type === 'text') {
        docPass.type = 'password';
        docIcon.className = 'fas fa-eye w3-text-theme';
    } else {
        docPass.type = 'text';
        docIcon.className = 'fas fa-eye-slash w3-text-theme';
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