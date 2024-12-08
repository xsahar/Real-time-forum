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

const MAX_RETRIES = 3;
const TIMEOUT = 5000;

export class APIService {
    constructor() {
        this.baseURL = '/api';
        this.token = localStorage.getItem('auth_token');
        console.log('[API] Service initialized, token:', this.token ? 'present' : 'missing');
    }

    async request(method, endpoint, data = null, retryCount = 0) {
        try {
            console.log(`[API] Making ${method} request to ${endpoint}`);
            
            const options = {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                signal: AbortSignal.timeout(TIMEOUT)
            };
    
            if (this.token) {
                options.headers['Authorization'] = `Bearer ${this.token}`;
            }
    
            if (data) {
                const formattedData = this.formatRequestData(endpoint, data);
                options.body = JSON.stringify(formattedData);
            }
    
            const response = await fetch(`${this.baseURL}/${endpoint}`, options);
            console.log(`[API] ${endpoint} status:`, response.status);
    
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
    
            const text = await response.text();
            
            // Handle empty response differently for auth/status endpoint
            if (endpoint === 'auth/status' && (!text || text.trim() === '')) {
                return { status: 'success' }; // Return success for empty auth check
            }
    
            if (!text) {
                throw new Error('Empty response');
            }
    
            const result = JSON.parse(text);
            
            if (result.status === 'error') {
                throw new Error(result.error || 'Request failed');
            }
    
            return result;
        } catch (error) {
            if (error.name === 'AbortError' && retryCount < MAX_RETRIES) {
                return this.request(method, endpoint, data, retryCount + 1);
            }
            throw error;
        }
    }

    formatRequestData(endpoint, data) {
        switch (endpoint) {
            case 'login':
                return this.formatLoginData(data);
            case 'register':
                return this.formatRegisterData(data);
            default:
                return data;
        }
    }

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

        // Validate register data
        const requiredFields = ['nickname', 'email', 'password', 'firstName', 'lastName', 'gender', 'age'];
        const missingFields = requiredFields.filter(field => !formattedData[field]);
        
        if (missingFields.length > 0) {
            throw new Error(`Missing required fields: ${missingFields.join(', ')}`);
        }

        if (formattedData.password.length < 6) {
            throw new Error('Password must be at least 6 characters');
        }
        if (!formattedData.email.includes('@')) {
            throw new Error('Invalid email format');
        }
        if (formattedData.nickname.length < 3) {
            throw new Error('Nickname must be at least 3 characters');
        }

        return formattedData;
    }

    async get(endpoint) {
        return this.request('GET', endpoint);
    }

            async post(endpoint, data) {
            return this.request('POST', endpoint, data);
        }
    }