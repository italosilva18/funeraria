from django.urls import path
from . import views

app_name = 'sistema_funeraria'


urlpatterns = [
    
    #path('', views.LoginView, name='login'),
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
    path('funcionarios/criar/', views.criar_funcionario, name='criar_funcionario'),
    path('funcionarios/editar/<int:id>/', views.editar_funcionario, name='editar_funcionarios'),
    path('funcionarios/excluir/<int:id>/', views.excluir_funcionario, name='excluir_funcionarios'),
    #Usuario
    #Planos
    path('plano/criar/', views.criar_plano, name='criar_plano'),
    path('plano/editar/<int:id>/', views.editar_plano, name='editar_plano'),
    path('plano/excluir/<int:id>/', views.excluir_plano, name='excluir_plano'),

]


handler404 = 'sistema_funeraria.views.error_404'
