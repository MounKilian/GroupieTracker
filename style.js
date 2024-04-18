const onglets = document.querySelectorAll('.onglet')
const contents = document.querySelectorAll('.content')
let index = 0

onglets.forEach(onglet => {
    onglet.addEventListener('click', () => {
        if (onglet.classList.contains('active')){
            return;
        } else {
            onglet.classList.add('active');
        }
        index = onglet.getAttribute('data-anim');
        for(let i = 0; i < onglets.length; i++) {
            if (onglets[i].getAttribute('data-anim') !== index) {
                onglets[i].classList.remove('active');
            }
        }
        for(let j = 0; j < contents.length; j++) {
            if (contents[j].getAttribute('data-anim') === index) {
                contents[j].classList.add('activeContent');
            } else {
                contents[j].classList.remove('activeContent');
            }
        }
    })
})