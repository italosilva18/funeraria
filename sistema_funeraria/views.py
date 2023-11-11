from django.shortcuts import redirect
from django.shortcuts import render
from .models import Cliente, Funcionario, Usuario, PlanoFunerario, ProdutoServico, Orcamento, Dashboard
from django.contrib.auth import authenticate, login
from django.contrib.auth.forms import AuthenticationForm
from django.contrib.auth.decorators import login_required



@login_required
def dashboard_view(request):
    clientes = Cliente.objects.count()
    funcionarios = Funcionario.objects.count()
    usuarios = Usuario.objects.count()
    planos = PlanoFunerario.objects.count()
    produtos_servicos = ProdutoServico.objects.count()
    orcamentos = Orcamento.objects.count()

    context = {
        'clientes': clientes,
        'funcionarios': funcionarios,
        'usuarios': usuarios,
        'planos': planos,
        'produtos_servicos': produtos_servicos,
        'orcamentos': orcamentos,
    }

    return render(request, 'sistema_funeraria/dashboard.html', context)

def login_view(request):
    if request.method == 'POST':
        form = AuthenticationForm(request, data=request.POST)
        if form.is_valid():
            username = form.cleaned_data.get('username')
            password = form.cleaned_data.get('password')
            user = authenticate(username=username, password=password)
            if user is not None:
                login(request, user)
                return redirect('dashboard/')
    else:
        form = AuthenticationForm()
    return render(request, 'sistema_funeraria/login.html', {'form': form})

@login_required
def lista_clientes(request):
    clientes = Cliente.objects.all()
    return render(request, 'sistema_funeraria/lista_clientes.html', {'clientes': clientes})
@login_required
def lista_funcionarios(request):
    funcionarios = Funcionario.objects.all()
    return render(request, 'sistema_funeraria/lista_funcionarios.html', {'funcionarios': funcionarios})
@login_required
def lista_usuarios(request):
    usuarios = Usuario.objects.all()
    return render(request, 'sistema_funeraria/lista_usuarios.html', {'usuarios': usuarios})
@login_required
def lista_planos(request):
    planos = PlanoFunerario.objects.all()
    return render(request, 'sistema_funeraria/lista_planos.html', {'planos': planos})
@login_required
def lista_produtos_servicos(request):
    produtos_servicos = ProdutoServico.objects.all()
    return render(request, 'sistema_funeraria/lista_produtos_servicos.html', {'produtos_servicos': produtos_servicos})
@login_required
def lista_orcamentos(request):
    orcamentos = Orcamento.objects.all()
    return render(request, 'sistema_funeraria/lista_orcamentos.html', {'orcamentos': orcamentos})
