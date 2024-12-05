// static/js/services/api.js
export class APIService {
    constructor() {
        this.baseURL = '/api';
        console.log('API Service initialized');
    }

    async request(method, endpoint, data = null) {
        try {
            // Validate request parameters
            if (!method || !endpoint) {
                throw new Error('Method and endpoint are required');
            }
    
            const token = localStorage.getItem('auth_token');
            const options = {
                method,
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    ...(token && { 'Authorization': `Bearer ${token}` })
                },
                // Add timeout
                signal: AbortSignal.timeout(5000) // 5 second timeout
            };
    
            if (data) {
                // Format data based on endpoint
                const formattedData = endpoint === 'register' 
                    ? await this.formatRegisterData(data)
                    : endpoint === 'login' 
                        ? await this.formatLoginData(data)
                        : data;
    
                // Debug log formatted request data
                console.log(`[API] ${endpoint} request:`, {
                    method,
                    endpoint,
                    data: {
                        ...formattedData,
                        password: formattedData.password ? '[PRESENT]' : '[MISSING]'
                    }
                });
    
                options.body = JSON.stringify(formattedData);
            }
    
            const response = await fetch(`${this.baseURL}/${endpoint}`, options);
            const result = await response.json();
    
            // Log response details
            console.log(`[API] ${endpoint} response:`, {
                status: response.status,
                ok: response.ok,
                data: result
            });
    
            // Handle different error cases
            if (!response.ok) {
                const error = new Error(result.error || `Request failed: ${response.statusText}`);
                error.status = response.status;
                error.data = result;
                throw error;
            }
    
            return result;
        } catch (error) {
            // Handle different error types
            if (error.name === 'AbortError') {
                console.error(`[API] ${endpoint} timeout`);
                throw new Error('Request timeout');
            }
            if (error.name === 'TypeError') {
                console.error(`[API] Network error:`, error);
                throw new Error('Network error - please check your connection');
            }
    
            console.error(`[API] ${endpoint} error:`, {
                message: error.message,
                status: error.status,
                data: error.data
            });
            throw error;
        }
    }

    formatRegisterData(data) {
        const formattedData = {
            nickname: data.nickname?.trim(),
            email: data.email?.trim(),
            password: data.password,
            firstName: data.firstName?.trim(),
            lastName: data.lastName?.trim(),
            gender: data.gender,
            age: parseInt(data.age)
        };

        this.validateRegisterData(formattedData);
        return formattedData;
    }

    formatLoginData(data) {
        return {
            identifier: data.identifier?.trim(),
            password: data.password
        };
    }

    validateRegisterData(data) {
        const requiredFields = ['nickname', 'email', 'password', 'firstName', 'lastName', 'gender', 'age'];
        const missingFields = requiredFields.filter(field => !data[field]);
        
        if (missingFields.length > 0) {
            throw new Error(`Missing required fields: ${missingFields.join(', ')}`);
        }

        if (data.password.length < 6) {
            throw new Error('Password must be at least 6 characters');
        }
        if (!data.email.includes('@')) {
            throw new Error('Invalid email format');
        }
        if (data.nickname.length < 3) {
            throw new Error('Nickname must be at least 3 characters');
        }
    }

    async get(endpoint) {
        return this.request('GET', endpoint);
    }

    async post(endpoint, data) {
        return this.request('POST', endpoint, data);
    }
}