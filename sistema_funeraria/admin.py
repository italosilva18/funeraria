from django.contrib import admin
from .models import Endereco, Contato, Cliente, Funcionario, Usuario, ProdutoServico, PlanoFunerario, ProdutoServicoPlano, Orcamento, ProdutoServicoOrcamento

# Personalize a interface de administração para o modelo Cliente
class ClienteAdmin(admin.ModelAdmin):
    list_display = ('user', 'endereco', 'contato')

# Personalize a interface de administração para o modelo Funcionario
class FuncionarioAdmin(admin.ModelAdmin):
    list_display = ('user', 'endereco', 'contato', 'cargo', 'setor', 'data_admissao')

# Personalize a interface de administração para o modelo Usuario
class UsuarioAdmin(admin.ModelAdmin):
    list_display = ('user', 'contato', 'nivel_acesso')

# Personalize a interface de administração para o modelo ProdutoServico
class ProdutoServicoAdmin(admin.ModelAdmin):
    list_display = ('nome', 'descricao', 'valor_unitario', 'tipo')

# Personalize a interface de administração para o modelo PlanoFunerario
class PlanoFunerarioAdmin(admin.ModelAdmin):
    list_display = ('nome', 'descricao', 'valor')

# Personalize a interface de administração para o modelo Orcamento
class OrcamentoAdmin(admin.ModelAdmin):
    list_display = ('cliente', 'valor_total')

# Registre seus modelos aqui
admin.site.register(Endereco)
admin.site.register(Contato)
admin.site.register(Cliente, ClienteAdmin)
admin.site.register(Funcionario, FuncionarioAdmin)
admin.site.register(Usuario, UsuarioAdmin)
admin.site.register(ProdutoServico, ProdutoServicoAdmin)
admin.site.register(PlanoFunerario, PlanoFunerarioAdmin)
admin.site.register(ProdutoServicoPlano)
admin.site.register(Orcamento, OrcamentoAdmin)
admin.site.register(ProdutoServicoOrcamento)
