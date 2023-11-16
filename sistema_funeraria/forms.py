from django.contrib.auth.models import User
from django import forms
from .models import Cliente, Contato, Funcionario, Usuario, PlanoFunerario, ProdutoServico, Orcamento
from django.contrib.auth.forms import UserCreationForm
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



class UserForm(UserCreationForm):
    email = forms.EmailField(required=True)

    class Meta:
        model = User
        fields = ("username", "email", "password1", "password2")

    def save(self, commit=True):
        user = super(UserForm, self).save(commit=False)
        user.email = self.cleaned_data['email']
        if commit:
            user.save()
        return user
