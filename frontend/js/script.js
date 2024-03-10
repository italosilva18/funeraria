// Exemplo de funcionalidade para exibir uma mensagem de confirmação ao clicar em um botão de exclusão
document.addEventListener('DOMContentLoaded', function() {
    const deleteButtons = document.querySelectorAll('.btn-delete');

    deleteButtons.forEach(button => {
        button.addEventListener('click', function() {
            const confirmation = confirm('Tem certeza que deseja excluir este item?');

            if (confirmation) {
                // Aqui você pode adicionar a lógica para excluir o item
                alert('Item excluído com sucesso!');
            } else {
                // Cancelou a exclusão
                alert('Exclusão cancelada.');
            }
        });
    });
});
