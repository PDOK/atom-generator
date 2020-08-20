function getXML () {
return fetch(window.location.href)
    .then(response => response.text())
    .then(data => {return data})
}

window.addEventListener('load', function () {
    getXML().then(function (data) {
    document.getElementById("atom-xml").innerHTML = data.replace(/\<\?xml-stylesheet.*?\>\n/, '').replace(/</g, "&lt;").replace(/>/g, "&gt;")
    Prism.highlightAll()
    let btn = document.getElementById("show")
    let wrapper = document.getElementById("xml-wrapper")
    // event listener for the show/hide atom xml button
    btn.onclick = function(){
        if (wrapper.style.display == "none"){
            wrapper.style.display = "block"
            btn.innerText = "Hide"
        }else{
            wrapper.style.display = "none"
            btn.innerText = "Show"
        }
    };
    })
});
