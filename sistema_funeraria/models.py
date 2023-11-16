from django.db import models
from django.contrib.auth.models import User

class Endereco(models.Model):
    rua = models.CharField(max_length=255)
    cidade = models.CharField(max_length=255)
    estado = models.CharField(max_length=255)
    cep = models.CharField(max_length=10)

    def __str__(self):
        return f'{self.rua}, {self.cidade}, {self.estado}, {self.cep}'

class Contato(models.Model):
    telefone = models.CharField(max_length=20)
    email = models.EmailField()

    def __str__(self):
        return f'{self.telefone}, {self.email}'

class Cliente(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    endereco = models.OneToOneField(Endereco, on_delete=models.CASCADE)
    contato = models.OneToOneField(Contato, on_delete=models.CASCADE)

    def __str__(self):
        return f'{self.user.username}'

class Dashboard(models.Model):
    title = models.CharField(max_length=200)
    description = models.TextField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return self.title

class Funcionario(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    endereco = models.OneToOneField(Endereco, on_delete=models.CASCADE)
    contato = models.OneToOneField(Contato, on_delete=models.CASCADE)
    cargo = models.CharField(max_length=255)
    setor = models.CharField(max_length=255)
    data_admissao = models.DateField()

    def __str__(self):
        return f'{self.user.username}, {self.cargo}'

class Usuario(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    contato = models.OneToOneField(Contato, on_delete=models.CASCADE)
    nivel_acesso = models.IntegerField()

    def __str__(self):
        return f'{self.user.username}, Nível de acesso: {self.nivel_acesso}'

class ProdutoServico(models.Model):
    TIPO_CHOICES = [
        ('P', 'Produto'),
        ('S', 'Serviço'),
    ]
    nome = models.CharField(max_length=255)
    descricao = models.TextField()
    valor_unitario = models.DecimalField(max_digits=6, decimal_places=2)
    tipo = models.CharField(max_length=1, choices=TIPO_CHOICES)

    def __str__(self):
        return f'{self.nome}, {self.get_tipo_display()}'

class PlanoFunerario(models.Model):
    nome = models.CharField(max_length=255, default='Nome padrão')
    descricao = models.TextField()
    valor = models.DecimalField(max_digits=6, decimal_places=2)
    produtos_servicos = models.ManyToManyField(ProdutoServico, through='ProdutoServicoPlano')

    def __str__(self):
        return self.nome

class ProdutoServicoPlano(models.Model):
    plano = models.ForeignKey(PlanoFunerario, on_delete=models.CASCADE)
    produto_servico = models.ForeignKey(ProdutoServico, on_delete=models.CASCADE)

    def __str__(self):
        return f'{self.plano.nome}, {self.produto_servico.nome}'

class Orcamento(models.Model):
    cliente = models.ForeignKey(Cliente, on_delete=models.CASCADE)
    produtos_servicos = models.ManyToManyField(ProdutoServico, through='ProdutoServicoOrcamento')
    valor_total = models.DecimalField(max_digits=6, decimal_places=2)

    def __str__(self):
        return f'Orçamento {self.id} - Cliente: {self.cliente.user.username}'

class ProdutoServicoOrcamento(models.Model):
    orcamento = models.ForeignKey(Orcamento, on_delete=models.CASCADE)
    produto_servico = models.ForeignKey(ProdutoServico, on_delete=models.CASCADE)

    def __str__(self):
        return f'{self.orcamento.id}, {self.produto_servico.nome}'
