/* ========================================
   Football Stats — App Module
   Toast, форматирование, мобильное меню
   ======================================== */

// ─────── Toast-уведомления ───────
function showToast(message, type = 'success') {
    const icons = { success: 'check-circle', error: 'exclamation-circle', info: 'info-circle' };
    const colors = { success: '#10b981', error: '#ef4444', info: '#3b82f6' };
    const toast = document.createElement('div');
    toast.className = 'toast';
    toast.innerHTML = `<i class="fas fa-${icons[type] || icons.info}"></i><span>${message}</span>`;
    toast.style.cssText = `
        position:fixed;top:24px;right:24px;background:${colors[type] || colors.info};
        color:#fff;padding:14px 24px;border-radius:12px;display:flex;align-items:center;
        gap:10px;font-weight:600;font-size:.95rem;box-shadow:0 8px 32px rgba(0,0,0,.4);
        z-index:9999;animation:slideIn .35s ease;font-family:'Inter',sans-serif;
    `;
    document.body.appendChild(toast);
    setTimeout(() => {
        toast.style.animation = 'slideOut .35s ease';
        setTimeout(() => toast.remove(), 350);
    }, 3500);
}

// ─────── Форматирование даты ───────
function formatDate(dateStr) {
    if (!dateStr) return '—';
    const d = new Date(dateStr);
    return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' });
}

function formatDateShort(dateStr) {
    if (!dateStr) return '—';
    const d = new Date(dateStr);
    return d.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' });
}

// ─────── Мобильное меню ───────
document.addEventListener('DOMContentLoaded', () => {
    const hamburger = document.querySelector('.hamburger');
    const navMenu = document.querySelector('.nav-menu');
    if (hamburger && navMenu) {
        hamburger.addEventListener('click', () => {
            hamburger.classList.toggle('active');
            navMenu.classList.toggle('active');
        });
    }
});

// ─────── CSS для toast-анимаций ───────
const toastStyle = document.createElement('style');
toastStyle.textContent = `
    @keyframes slideIn { from{transform:translateX(120%);opacity:0}to{transform:translateX(0);opacity:1} }
    @keyframes slideOut { from{transform:translateX(0);opacity:1}to{transform:translateX(120%);opacity:0} }
`;
document.head.appendChild(toastStyle);
