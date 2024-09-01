document.addEventListener('DOMContentLoaded', function() {
    const updateProductButtons = document.querySelectorAll('.updateProductButton');
    const drawer = document.getElementById('drawer');
    const updateForm = document.getElementById('updateForm');
    const overlay = document.getElementById('overlay');
    const updateMonthlyBudgetButton = document.getElementById('updateMonthlyBudget');

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

    if (updateMonthlyBudgetButton) {
        updateMonthlyBudgetButton.addEventListener('click', () => {
            drawer.classList.add('open');
            overlay.classList.add('visible');

            document.getElementById('category').textContent = '📅 Monthly'; 
        });
    }

    overlay.addEventListener('click', () => {
        drawer.classList.remove('open');
        overlay.classList.remove('visible');
    });

    updateForm.addEventListener('submit', function(event) {
        event.preventDefault();

        const categoryWithIcon = document.getElementById('category').textContent;
        const category = categoryWithIcon.slice(3).toLowerCase();
        const budgetAmount = updateForm.budgetAmount.value;
        const data = {
            category: category,
            budgetAmount: budgetAmount,
        };

        fetch('/update-budget', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(result => {
            window.location.reload();
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
});
