from django import forms
from .models import Cliente, Contato, Funcionario, PlanoFunerario

class PlanoForm(forms.ModelForm):
    class Meta:
        model = PlanoFunerario
        fields = ['nome', 'descricao', 'valor', 'produtos_servicos']

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
