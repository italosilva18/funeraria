from django.urls import path
from . import views

app_name = 'sistema_funeraria'


urlpatterns = [
    
    #path('', views.LoginView.as_view, name='login'),
    path('login/', views.LoginView.as_view(), name='login'),
    path('logout/', views.logout_view, name='logout'),
    path('dashboard/', views.dashboard_view, name='dashboard'),
    path('lista_clientes/', views.ListaClientesView.as_view(), name='lista_clientes'),
    path('lista_funcionarios/', views.ListaFuncionariosView.as_view(), name='lista_funcionarios'),
    path('lista_usuarios/', views.lista_usuariosView.as_view(), name='lista_usuarios'),
    path('lista_planos/', views.lista_planosView.as_view(), name='lista_planos'),
    path('lista_produtos_servicos/', views.lista_produtos_servicosView.as_view(), name='lista_produtos_servicos'),
    path('lista_orcamentos/', views.lista_orcamentosView.as_view(), name='lista_orcamentos'),
    #Clientes
    path('clientes/criar/', views.criar_cliente, name='criar_cliente'),
    path('clientes/editar/<int:id>/', views.editar_cliente, name='editar_cliente'),
    path('clientes/excluir/<int:id>/', views.excluir_cliente, name='excluir_cliente'),
    #Funcionario
    path('funcionario/criar/', views.criar_funcionario, name='criar_funcionario'),
    path('funcionarios/editar/<int:id>/', views.editar_funcionario, name='editar_funcionario'),
    path('funcionarios/excluir/<int:id>/', views.excluir_funcionario, name='excluir_funcionario'),

    #Usuario
    path('usuarios/criar/', views.criar_usuarios, name='criar_usuarios'),
    path('usuarios/editar/<int:id>/', views.editar_usuarios, name='editar_usuarios'),
    path('usuarios/excluir/<int:id>/', views.excluir_usuarios, name='excluir_usuarios'),
    #Planos
    path('plano/criar/', views.criar_plano, name='criar_plano'),
    path('plano/editar/<int:id>/', views.editar_plano, name='editar_plano'),
    path('plano/excluir/<int:id>/', views.excluir_plano, name='excluir_plano'),
    #orcamentos
    path('orcamento/criar/', views.criar_orcamento, name='criar_orcamento'),
    path('orcamento/editar/<int:id>/', views.editar_orcamento, name='editar_orcamento'),
    path('orcamento/excluir/<int:id>/', views.excluir_orcamento, name='excluir_orcamento'),
    #Produto_Servico
    path('produto_servico/criar/', views.criar_produto_servico, name='criar_produto_servico'),
    path('produto_servico/editar/<int:id>/', views.editar_produto_servico, name='editar_produto_servico'),
    path('produto_servico/excluir/<int:id>/', views.excluir_produto_servico, name='excluir_produto_servico'),
]


handler404 = 'sistema_funeraria.views.error_404'
