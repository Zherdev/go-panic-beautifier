function  focus() {
    document.getElementById('panic-log').focus();
    window.setTimeout(focus, 200);
}

function copy() {
    /* Get the text field */
    let toCopy = document.getElementById("result").textContent;

    const tmp = document.createElement('textarea');
    tmp.value = toCopy;

    document.body.appendChild(tmp);
    tmp.select();
    document.execCommand('copy');
    document.body.removeChild(tmp);
}

window.addEventListener("load", function(){
    window.setTimeout(function () {
        document.getElementById('panic-log').focus();
    }, 200);

    let logForm = document.getElementById("panic-input-form");

    logForm.onpaste = function () {
        setTimeout(() => {logForm.submit();}, 200);
    };
});
