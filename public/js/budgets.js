document.addEventListener('DOMContentLoaded', function() {
    const updateProductButtons = document.querySelectorAll('.updateProductButton');
    const drawer = document.getElementById('drawer');
    const updateForm = document.getElementById('updateForm');
    const overlay = document.getElementById('overlay');

    updateProductButtons.forEach(button => {
        button.addEventListener('click', () => {
            drawer.classList.add('open');
            overlay.classList.add('visible');

            const row = button.closest('tr');
            const category = row.querySelector('td:nth-child(1) .text-base').textContent;
            const budgetAmount = row.querySelector('td:nth-child(2)').textContent.trim().replace('$', '');

            document.getElementById('category').textContent = category;
            updateForm.budgetAmount.value = budgetAmount;
        });
    });

    overlay.addEventListener('click', () => {
        drawer.classList.remove('open');
        overlay.classList.remove('visible');
    });

    updateForm.addEventListener('submit', function(event) {
        event.preventDefault();

        const category = document.getElementById('category').textContent;
        const budgetAmount = updateForm.budgetAmount.value;

        const rows = document.querySelectorAll('tbody tr');
        rows.forEach(row => {
            if (row.querySelector('td:nth-child(1) .text-base').textContent === category) {
                row.querySelector('td:nth-child(2)').textContent = `$${budgetAmount}`;
            }
        });

        drawer.classList.remove('open');
        overlay.classList.remove('visible');
    });
});