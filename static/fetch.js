document.addEventListener('DOMContentLoaded', () => {
    fetch("/api/categories")
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            const categorySelect = document.getElementById("category-select");
            if (categorySelect) {
                categorySelect.innerHTML = '<option value="">All Categories</option>'; // Reset options
                data.categories.forEach(category => {
                    const option = document.createElement("option");
                    option.value = category.id;
                    option.textContent = category.name;
                    categorySelect.appendChild(option);
                });
            }
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
});