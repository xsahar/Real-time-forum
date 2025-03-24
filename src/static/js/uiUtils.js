export function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        return new Promise((resolve) => {
            const later = () => {
                clearTimeout(timeout);
                resolve(func(...args));
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        });
    };
}
