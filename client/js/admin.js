const form = document.forms[0];
form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const formData = new FormData(e.target);

    try {
        const response = await fetch('http://localhost:8073/api/add-product', {
            method: 'POST',
            body: formData
        })

        if (response.ok) {
            const result = await response.json();
            console.log('Success:', result);
        } else {
            console.error('Server error:', response.statusText);
        }
    } catch(error) {
        console.error('try error:', error);
    }
});