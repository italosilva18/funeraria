from django.urls import path
from . import views

app_name = 'sistema_funeraria'

urlpatterns = [
    
    path('', views.login_view, name='login'),
    path('dashboard/', views.dashboard_view, name='dashboard'),
    path('lista_clientes/', views.lista_clientes, name='lista_clientes'),
    path('lista_funcionarios/', views.lista_funcionarios, name='lista_funcionarios'),
    path('lista_usuarios/', views.lista_usuarios, name='lista_usuarios'),
    path('lista_planos/', views.lista_planos, name='lista_planos'),
    path('lista_produtos_servicos/', views.lista_produtos_servicos, name='lista_produtos_servicos'),
    path('lista_orcamentos/', views.lista_orcamentos, name='lista_orcamentos'),
]
