from django import forms
from .models import Cliente, Contato, Funcionario, Usuario, PlanoFunerario, ProdutoServico, Orcamento

class UsuarioForm(forms.ModelForm):
    class Meta:
        model = Usuario
        fields = ['user', 'contato', 'nivel_acesso']

class PlanoForm(forms.ModelForm):
    class Meta:
        model = PlanoFunerario
        fields = ['nome', 'descricao', 'valor', 'produtos_servicos']

class ProdutoServicoForm(forms.ModelForm):
    class Meta:
        model = ProdutoServico
        fields = ['nome', 'descricao', 'valor_unitario', 'tipo']

class OrcamentoForm(forms.ModelForm):
    class Meta:
        model = Orcamento
        fields = ['cliente', 'produtos_servicos', 'valor_total']

class ContatoForm(forms.ModelForm):
    class Meta:
        model = Contato
        fields = ['email', 'telefone']

class ClienteForm(forms.ModelForm):
    class Meta:
        model = Cliente
        fields = ['user', 'endereco', 'contato']


class FuncionarioForm(forms.ModelForm):
    class Meta:
        model = Funcionario
        fields = ['user', 'endereco', 'contato', 'cargo', 'setor', 'data_admissao']
