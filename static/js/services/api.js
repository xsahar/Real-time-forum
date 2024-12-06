// api.js
const ENDPOINTS = {
    auth: {
        login: 'login',
        register: 'register',
        status: 'auth/status'
    },
    users: {
        list: 'users',
        online: 'users/online'
    },
    posts: {
        list: 'posts',
        create: 'posts/create'
    }
};

export class APIService {
    constructor() {
        this.baseURL = '/api';
        this.token = localStorage.getItem('auth_token');
        console.log('API Service initialized, token:', this.token ? 'present' : 'missing');
    }

    async request(method, endpoint, data = null) {
        try {
            if (!method || !endpoint) {
                throw new Error('Method and endpoint are required');
            }

            // Get fresh token on each request
            const token = localStorage.getItem('auth_token');
            console.log(`[API] Token for ${endpoint}:`, token ? 'present' : 'missing');

            const options = {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    ...(token && { 'Authorization': `Bearer ${token}` })
                },
                signal: AbortSignal.timeout(5000)
            };

            if (data) {
                const formattedData = 
                    endpoint === 'register' ? this.formatRegisterData(data)
                    : endpoint === 'login' ? this.formatLoginData(data)
                    : data;

                console.log(`[API] ${endpoint} request:`, {
                    method,
                    endpoint,
                    hasToken: !!token,
                    data: {
                        ...formattedData,
                        password: '[REDACTED]'
                    }
                });

                options.body = JSON.stringify(formattedData);
            }

            const response = await fetch(`${this.baseURL}/${endpoint}`, options);
            const text = await response.text();
            
            let result;
            try {
                result = text ? JSON.parse(text) : {};
            } catch (parseError) {
                console.error(`[API] JSON parse error:`, parseError);
                throw new Error('Invalid response format');
            }

            // Handle auth token in response
            if (result?.data?.token) {
                localStorage.setItem('auth_token', result.data.token);
                console.log('[API] Token stored successfully');
            }

            if (!response.ok) {
                throw new Error(result.error || `Request failed: ${response.statusText}`);
            }

            return result;
        } catch (error) {
            console.error(`[API] ${endpoint} error:`, error);
            throw error;
        }
    }

    // Rest of methods remain the same
    formatLoginData(data) {
        if (!data.identifier || !data.password) {
            throw new Error('Identifier and password are required');
        }
        return {
            identifier: data.identifier.trim(),
            password: data.password
        };
    }

    formatRegisterData(data) {
        const formattedData = {
            nickname: data.nickname?.trim(),
            email: data.email?.trim().toLowerCase(),
            password: data.password,
            firstName: data.firstName?.trim(),
            lastName: data.lastName?.trim(),
            gender: data.gender,
            age: parseInt(data.age)
        };
        this.validateRegisterData(formattedData);
        return formattedData;
    }

    validateRegisterData(data) {
        const requiredFields = ['nickname', 'email', 'password', 'firstName', 'lastName', 'gender', 'age'];
        const missingFields = requiredFields.filter(field => !data[field]);
        if (missingFields.length > 0) {
            throw new Error(`Missing required fields: ${missingFields.join(', ')}`);
        }
        if (data.password.length < 6) throw new Error('Password must be at least 6 characters');
        if (!data.email.includes('@')) throw new Error('Invalid email format');
        if (data.nickname.length < 3) throw new Error('Nickname must be at least 3 characters');
    }

    async get(endpoint) {
        return this.request('GET', endpoint);
    }

    async post(endpoint, data) {
        return this.request('POST', endpoint, data);
    }
}