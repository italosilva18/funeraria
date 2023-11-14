from django.shortcuts import render, redirect, get_object_or_404
from django.views import View
from django.views.generic.list import ListView
from .models import Cliente, Funcionario, Usuario, PlanoFunerario, ProdutoServico, Orcamento
from django.contrib.auth import authenticate, login, logout
from django.contrib.auth.forms import AuthenticationForm
from django.contrib.auth.mixins import LoginRequiredMixin
from django.contrib.auth.decorators import login_required
from django.http import HttpResponseRedirect
from .forms import ClienteForm, ContatoForm, FuncionarioForm, PlanoForm, UsuarioForm


def error_404(request, exception):
    return redirect('sistema_funeraria:login')


@login_required
def dashboard_view(request):
    context = {
        'clientes': Cliente.objects.count(),
        'funcionarios': Funcionario.objects.count(),
        'usuarios': Usuario.objects.count(),
        'planos': PlanoFunerario.objects.count(),
        'produtos_servicos': ProdutoServico.objects.count(),
        'orcamentos': Orcamento.objects.count(),
    }

    return render(request, 'sistema_funeraria/dashboard.html', context)


class LoginView(View):
    def get(self, request):
        return render(request, 'sistema_funeraria/login.html', {'form': AuthenticationForm()})

    def post(self, request):
        form = AuthenticationForm(request, data=request.POST)
        if form.is_valid():
            username = form.cleaned_data.get('username')
            password = form.cleaned_data.get('password')
            user = authenticate(username=username, password=password)
            if user is not None and user.is_active:
                login(request, user)
                return redirect('sistema_funeraria:dashboard')
        return render(request, 'sistema_funeraria/login.html', {'form': form})


def logout_view(request):
    logout(request)
    return redirect('sistema_funeraria:login')



#clientes
class ListaClientesView(LoginRequiredMixin, ListView):
    model = Cliente
    template_name = 'sistema_funeraria/cliente/lista_clientes.html'
    context_object_name = 'clientes'

def criar_cliente(request):
    if request.method == "POST":
        form = ClienteForm(request.POST)
        if form.is_valid():
            cliente = form.save()
            return redirect('sistema_funeraria:lista_clientes')
    else:
        form = ClienteForm()
    return render(request, 'sistema_funeraria/cliente/criar_cliente.html', {'form': form})


def editar_cliente(request, id):
    cliente = get_object_or_404(Cliente, id=id)
    if request.method == "POST":
        form = ClienteForm(request.POST, instance=cliente)
        contato_form = ContatoForm(request.POST, instance=cliente.contato)
        if form.is_valid() and contato_form.is_valid():
            form.save()
            contato_form.save()
            return HttpResponseRedirect('/lista_clientes/')
    else:
        form = ClienteForm(instance=cliente)
        contato_form = ContatoForm(instance=cliente.contato)
    return render(request, 'sistema_funeraria/cliente/editar_cliente.html', {'form': form, 'contato_form': contato_form})

def excluir_cliente(request, id):
    cliente = get_object_or_404(Cliente, id=id)
    if request.method == "POST":
        cliente.delete()
        return HttpResponseRedirect('/lista_clientes/')
    return render(request, 'sistema_funeraria/cliente/excluir_cliente.html', {'cliente': cliente})



class ListaFuncionariosView(LoginRequiredMixin, ListView):
    model = Funcionario
    template_name = 'sistema_funeraria/funcionario/lista_funcionarios.html'
    context_object_name = 'funcionarios'

def criar_funcionario(request):
    if request.method == "POST":
        form = FuncionarioForm(request.POST)
        if form.is_valid():
            funcionario = form.save()
            return redirect('sistema_funeraria:lista_funcionarios')
    else:
        form = FuncionarioForm()
    return render(request, 'sistema_funeraria/funcionario/criar_funcionario.html', {'form': form})

def editar_funcionario(request, id):
    funcionario = get_object_or_404(Funcionario, id=id)
    if request.method == "POST":
        form = FuncionarioForm(request.POST, instance=funcionario)
        if form.is_valid():
            funcionario = form.save()
            return redirect('sistema_funeraria:lista_funcionarios')
    else:
        form = FuncionarioForm(instance=funcionario)
    return render(request, 'sistema_funeraria/funcionario/editar_funcionario.html', {'form': form})

def excluir_funcionario(request, id):
    funcionario = get_object_or_404(Funcionario, id=id)
    if request.method == "POST":
        funcionario.delete()
        return redirect('sistema_funeraria:lista_funcionarios')
    return render(request, 'sistema_funeraria/funcionario/excluir_funcionario.html', {'funcionario': funcionario})



#usuarios
class lista_usuariosView(LoginRequiredMixin, ListView):
    model = Usuario
    template_name = 'sistema_funeraria/usuarios/lista_usuarios.html'
    context_object_name = 'usuarios'

def criar_usuarios(request):
    if request.method == "POST":
        form = UsuarioForm(request.POST)
        if form.is_valid():
            usuario = form.save()
            return redirect('sistema_funeraria:lista_usuarios')
    else:
        form = UsuarioForm()
    return render(request, 'sistema_funeraria/usuarios/criar_usuario.html', {'form': form})

def editar_usuarios(request, id):
    usuario = get_object_or_404(Usuario, id=id)
    if request.method == "POST":
        form = UsuarioForm(request.POST, instance=usuario)
        if form.is_valid():
            usuario = form.save()
            return redirect('sistema_funeraria:lista_usuarios')
    else:
        form = UsuarioForm(instance=usuario)
    return render(request, 'sistema_funeraria/usuarios/editar_usuario.html', {'form': form})

def excluir_usuarios(request, id):
    usuario = get_object_or_404(Usuario, id=id)
    if request.method == "POST":
        usuario.delete()
        return redirect('sistema_funeraria:lista_usuarios')
    return render(request, 'sistema_funeraria/usuarios/excluir_usuario.html', {'usuario': usuario})


class lista_planosView(LoginRequiredMixin,ListView):
    model = PlanoFunerario
    template_name = 'sistema_funeraria/plano/lista_planos.html'
    context_object_name = 'planos'


def criar_plano(request):
    if request.method == "POST":
        form = PlanoForm(request.POST)
        if form.is_valid():
            plano = form.save()
            return redirect('sistema_funeraria:lista_planos')
    else:
        form = PlanoForm()
    return render(request, 'sistema_funeraria/plano/criar_plano.html', {'form': form})

def editar_plano(request, id):
    plano = get_object_or_404(PlanoFunerario, id=id)
    if request.method == "POST":
        form = PlanoForm(request.POST, instance=plano)
        if form.is_valid():
            plano = form.save()
            return redirect('sistema_funeraria:lista_planos')
    else:
        form = PlanoForm(instance=plano)
    return render(request, 'sistema_funeraria/plano/editar_plano.html', {'form': form})

def excluir_plano(request, id):
    plano = get_object_or_404(PlanoFunerario, id=id)
    if request.method == "POST":
        plano.delete()
        return redirect('sistema_funeraria:lista_planos')
    return render(request, 'sistema_funeraria/plano/excluir_plano.html', {'plano': plano})

class lista_produtos_servicosView(LoginRequiredMixin,ListView):
    model = ProdutoServico
    template_name = 'sistema_funeraria/produtos_servicos/lista_produtos_servicos.html'    
    context_object_name = 'produtos_servicos'
    
class lista_orcamentosView(LoginRequiredMixin, ListView):
    model = Orcamento
    template_name = 'sistema_funeraria/orcamentos/lista_orcamentos.html'
    context_object_name = 'orcamentos'


