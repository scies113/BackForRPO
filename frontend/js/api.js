/* ========================================
   Football Stats — API Module
   Все запросы к бэкенду
   ======================================== */

const API_BASE = '/api';

const API = {
    // ─────── Авторизация ───────
    async login(email, password) {
        return this._post('/login', { email, password });
    },

    async register(username, email, password) {
        return this._post('/register', { username, email, password });
    },

    async logout() {
        return this._post('/logout');
    },

    async getMe() {
        return this._get('/me');
    },

    // ─────── Матчи ───────
    async getMatches() {
        return this._get('/matches');
    },

    async getMatch(id) {
        return this._get(`/matches/${id}`);
    },

    async createMatch(data) {
        return this._post('/matches', data);
    },

    async updateMatch(id, data) {
        return this._put(`/matches/${id}`, data);
    },

    async deleteMatch(id) {
        return this._delete(`/matches/${id}`);
    },

    // ─────── Прогнозы ИИ ───────
    async getPrediction(matchId) {
        return this._post(`/predict/${matchId}`);
    },

    // ─────── Аудит ───────
    async getAuditLogs(params = {}) {
        const query = new URLSearchParams();
        if (params.user_id) query.set('user_id', params.user_id);
        if (params.action) query.set('action', params.action);
        if (params.date_from) query.set('date_from', params.date_from);
        if (params.date_to) query.set('date_to', params.date_to);
        if (params.page) query.set('page', params.page);
        if (params.limit) query.set('limit', params.limit);
        const qs = query.toString();
        return this._get('/audit' + (qs ? '?' + qs : ''));
    },

    // ─────── Внутренние методы ───────
    async _get(path) {
        return this._request('GET', path);
    },

    async _post(path, body) {
        return this._request('POST', path, body);
    },

    async _put(path, body) {
        return this._request('PUT', path, body);
    },

    async _delete(path) {
        return this._request('DELETE', path);
    },

    async _request(method, path, body) {
        const opts = {
            method,
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',   // отправляем куки
        };
        if (body) {
            opts.body = JSON.stringify(body);
        }
        try {
            const res = await fetch(API_BASE + path, opts);
            const data = await res.json().catch(() => null);
            if (!res.ok) {
                const msg = data?.message || data?.error || `Ошибка ${res.status}`;
                throw { status: res.status, message: msg, code: data?.code };
            }
            return data;
        } catch (err) {
            if (err.status) throw err;
            throw { status: 0, message: 'Сервер недоступен' };
        }
    }
};
