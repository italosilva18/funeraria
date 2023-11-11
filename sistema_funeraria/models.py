from django.db import models
from django.contrib.auth.models import User

class Endereco(models.Model):
    rua = models.CharField(max_length=255)
    cidade = models.CharField(max_length=255)
    estado = models.CharField(max_length=255)
    cep = models.CharField(max_length=10)

class Contato(models.Model):
    telefone = models.CharField(max_length=20)
    email = models.EmailField()

class Cliente(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    endereco = models.OneToOneField(Endereco, on_delete=models.CASCADE)
    contato = models.OneToOneField(Contato, on_delete=models.CASCADE)
    
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

class Usuario(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    contato = models.OneToOneField(Contato, on_delete=models.CASCADE)
    nivel_acesso = models.IntegerField()

class ProdutoServico(models.Model):
    TIPO_CHOICES = [
        ('P', 'Produto'),
        ('S', 'Serviço'),
    ]
    nome = models.CharField(max_length=255)
    descricao = models.TextField()
    valor_unitario = models.DecimalField(max_digits=6, decimal_places=2)
    tipo = models.CharField(max_length=1, choices=TIPO_CHOICES)

class PlanoFunerario(models.Model):
    nome = models.CharField(max_length=255, default='Nome padrão')
    descricao = models.TextField()
    valor = models.DecimalField(max_digits=6, decimal_places=2)
    produtos_servicos = models.ManyToManyField(ProdutoServico, through='ProdutoServicoPlano')

class ProdutoServicoPlano(models.Model):
    plano = models.ForeignKey(PlanoFunerario, on_delete=models.CASCADE)
    produto_servico = models.ForeignKey(ProdutoServico, on_delete=models.CASCADE)

class Orcamento(models.Model):
    cliente = models.ForeignKey(Cliente, on_delete=models.CASCADE)
    produtos_servicos = models.ManyToManyField(ProdutoServico, through='ProdutoServicoOrcamento')
    valor_total = models.DecimalField(max_digits=6, decimal_places=2)

class ProdutoServicoOrcamento(models.Model):
    orcamento = models.ForeignKey(Orcamento, on_delete=models.CASCADE)
    produto_servico = models.ForeignKey(ProdutoServico, on_delete=models.CASCADE)
