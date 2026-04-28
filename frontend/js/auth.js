/* ========================================
   Football Stats — Auth Module
   Проверка авторизации, RBAC, редирект, UI
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
            const roleBadge = { admin: '👑', operator: '⚙️', analyst: '📊', user: '👤' };
            navActions.innerHTML =
                '<span class="nav-user">' +
                    '<span class="nav-role">' + (roleBadge[this.user.role] || '👤') + '</span>' +
                    '<span class="nav-username">' + this.user.username + '</span>' +
                '</span>' +
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
