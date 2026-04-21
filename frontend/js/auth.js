/* ========================================
   Football Stats — Auth Module
   Проверка авторизации, редирект, UI
   ======================================== */

const Auth = {
    user: null,

    // Проверить авторизацию, загрузить данные пользователя
    async check() {
        try {
            this.user = await API.getMe();
            this.updateNav();
            return true;
        } catch {
            this.user = null;
            return false;
        }
    },

    // Если не авторизован — редирект на логин
    async requireAuth() {
        const ok = await this.check();
        if (!ok) {
            window.location.href = '/login';
            return false;
        }
        return true;
    },

    // Если не admin — редирект
    async requireAdmin() {
        const ok = await this.requireAuth();
        if (!ok) return false;
        if (this.user.role !== 'admin') {
            showToast('Доступ только для администратора', 'error');
            window.location.href = '/dashboard';
            return false;
        }
        return true;
    },

    // Проверка роли (admin или operator)
    canEdit() {
        return this.user && (this.user.role === 'admin' || this.user.role === 'operator');
    },

    canPredict() {
        return this.user && ['admin', 'operator', 'analyst'].includes(this.user.role);
    },

    // Обновить навигацию — показать имя пользователя
    updateNav() {
        const navActions = document.querySelector('.nav-actions');
        if (!navActions) return;

        if (this.user) {
            const roleBadge = {
                admin: '👑', operator: '⚙️', analyst: '📊', user: '👤'
            };
            navActions.innerHTML = `
                <span class="nav-user">
                    <span class="nav-role">${roleBadge[this.user.role] || '👤'}</span>
                    <span class="nav-username">${this.user.username}</span>
                </span>
                <button class="btn btn-outline btn-sm" onclick="Auth.logout()">Выйти</button>
            `;
        } else {
            navActions.innerHTML = `
                <a href="/login" class="btn btn-outline">Войти</a>
                <a href="/login#register" class="btn btn-primary">Регистрация</a>
            `;
        }
    },

    async logout() {
        try {
            await API.logout();
        } catch {}
        this.user = null;
        showToast('Вы вышли из системы', 'success');
        setTimeout(() => window.location.href = '/', 500);
    }
};
