const modal = document.getElementById('modal');
const modalTitle = document.getElementById('modalTitle');
const modalBody = document.getElementById('modalBody');

export const showModal = (title, content) => {
    if (modalTitle && modalBody && modal) {
        modalTitle.textContent = title;
        modalBody.innerHTML = content;
        modal.style.display = 'flex';
    }
};

export const closeModal = () => {
    if (modal) {
        modal.style.display = 'none';
    }
};

window.onclick = function (event) {
    if (event.target === modal) {
        closeModal();
    }
};

window.closeModal = closeModal;
