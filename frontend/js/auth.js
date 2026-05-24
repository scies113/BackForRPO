/* ========================================
   Football Stats — Auth Module
   Проверка авторизации, RBAC, редирект, UI
   ======================================== */

const Auth = {
    user: null,

    // SVG иконки для ролей (профессиональные, заменяют emoji)
    _roleIcons: {
        admin: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#00e676" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg>',
        operator: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#00e676" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>',
        analyst: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#00e676" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>',
        user: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="#00e676" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>'
    },

    // Проверить авторизацию, загрузить данные пользователя
    async check() {
        try {
            this.user = await API.getMe();
            this.updateNav();
            return true;
        } catch {
            this.user = null;
            this.updateNav();
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

    // ─── RBAC: Универсальный route guard ───
    // Проверяет, что текущий пользователь имеет одну из разрешённых ролей.
    // Если роль не совпадает — редирект на /dashboard с уведомлением.
    async requireRole(...roles) {
        const ok = await this.requireAuth();
        if (!ok) return false;

        if (!roles.includes(this.user.role)) {
            showToast('Недостаточно прав для доступа к этой странице', 'error');
            window.location.href = '/dashboard';
            return false;
        }
        return true;
    },

    // Обратная совместимость: requireAdmin()
    async requireAdmin() {
        return this.requireRole('admin');
    },

    // ─── Матрица доступов ───

    // CRUD матчей: admin, operator
    canEdit() {
        if (!this.user) return false;
        return ['admin', 'operator'].includes(this.user.role);
    },

    // Генерация прогнозов: admin, operator, analyst
    // user может только ПРОСМАТРИВАТЬ уже созданные прогнозы
    canPredict() {
        if (!this.user) return false;
        return ['admin', 'operator', 'analyst'].includes(this.user.role);
    },

    // ─── Обновление навигации (Header) ───
    updateNav() {
        const navActions = document.querySelector('.nav-actions');

        // Очищаем старые классы ролей у body
        document.body.classList.remove('role-admin', 'role-operator', 'role-analyst', 'role-user', 'role-guest');

        const userRole = this.user ? this.user.role : 'guest';
        document.body.classList.add('role-' + userRole);

        // ─── Nav-actions: кнопки Войти / Профиль+Выйти ───
        if (this.user && navActions) {
            const roleIcon = this._roleIcons[this.user.role] || this._roleIcons.user;
            navActions.innerHTML =
                '<a href="/dashboard" class="nav-user-link">' +
                    '<span class="nav-user">' +
                        '<span class="nav-role">' + roleIcon + '</span>' +
                        '<span class="nav-username">' + this.user.username + '</span>' +
                    '</span>' +
                '</a>' +
                '<button class="btn btn-outline btn-sm" onclick="Auth.logout()">Выйти</button>';
        } else if (navActions) {
            navActions.innerHTML =
                '<a href="/login" class="btn btn-outline">Войти</a>' +
                '<a href="/login#register" class="btn btn-primary">Регистрация</a>';
        }

        // ─── RBAC: Показ/скрытие элементов с data-role ───
        // Вместо CSS [data-role]{display:none!important} (который JS не может переопределить),
        // используем класс .role-hidden который добавляем/убираем через JS.
        document.querySelectorAll('[data-role]').forEach(function(el) {
            var allowed = el.getAttribute('data-role').split(',');
            if (allowed.indexOf(userRole) !== -1) {
                // Роль разрешена — убираем скрытие
                el.classList.remove('role-hidden');
            } else {
                // Роль не разрешена — скрываем
                el.classList.add('role-hidden');
            }
        });
    },

    async logout() {
        try { await API.logout(); } catch(e) {}
        this.user = null;
        showToast('Вы вышли из системы', 'success');
        setTimeout(function() { window.location.href = '/'; }, 500);
    }
};

// ─── Начальное скрытие: сразу скрыть все data-role элементы до проверки auth ───
// Это выполняется синхронно при загрузке скрипта, ДО вызова Auth.check()
(function() {
    document.querySelectorAll('[data-role]').forEach(function(el) {
        el.classList.add('role-hidden');
    });
})();
