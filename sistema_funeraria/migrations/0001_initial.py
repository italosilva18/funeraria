# Generated by Django 4.2.7 on 2023-11-10 01:32

from django.conf import settings
from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
    ]

    operations = [
        migrations.CreateModel(
            name='Cliente',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
            ],
        ),
        migrations.CreateModel(
            name='Contato',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('telefone', models.CharField(max_length=20)),
                ('email', models.EmailField(max_length=254)),
            ],
        ),
        migrations.CreateModel(
            name='Endereco',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('rua', models.CharField(max_length=255)),
                ('cidade', models.CharField(max_length=255)),
                ('estado', models.CharField(max_length=255)),
                ('cep', models.CharField(max_length=10)),
            ],
        ),
        migrations.CreateModel(
            name='Orcamento',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('valor_total', models.DecimalField(decimal_places=2, max_digits=6)),
                ('cliente', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.cliente')),
            ],
        ),
        migrations.CreateModel(
            name='PlanoFunerario',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('descricao', models.TextField()),
                ('valor', models.DecimalField(decimal_places=2, max_digits=6)),
            ],
        ),
        migrations.CreateModel(
            name='ProdutoServico',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('nome', models.CharField(max_length=255)),
                ('descricao', models.TextField()),
                ('valor_unitario', models.DecimalField(decimal_places=2, max_digits=6)),
                ('tipo', models.CharField(choices=[('P', 'Produto'), ('S', 'Serviço')], max_length=1)),
            ],
        ),
        migrations.CreateModel(
            name='Usuario',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('nivel_acesso', models.IntegerField()),
                ('contato', models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.contato')),
                ('user', models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.CreateModel(
            name='ProdutoServicoPlano',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('plano', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.planofunerario')),
                ('produto_servico', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.produtoservico')),
            ],
        ),
        migrations.CreateModel(
            name='ProdutoServicoOrcamento',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('orcamento', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.orcamento')),
                ('produto_servico', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.produtoservico')),
            ],
        ),
        migrations.AddField(
            model_name='planofunerario',
            name='produtos_servicos',
            field=models.ManyToManyField(through='sistema_funeraria.ProdutoServicoPlano', to='sistema_funeraria.produtoservico'),
        ),
        migrations.AddField(
            model_name='orcamento',
            name='produtos_servicos',
            field=models.ManyToManyField(through='sistema_funeraria.ProdutoServicoOrcamento', to='sistema_funeraria.produtoservico'),
        ),
        migrations.CreateModel(
            name='Funcionario',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('cargo', models.CharField(max_length=255)),
                ('setor', models.CharField(max_length=255)),
                ('data_admissao', models.DateField()),
                ('contato', models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.contato')),
                ('endereco', models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.endereco')),
                ('user', models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.AddField(
            model_name='cliente',
            name='contato',
            field=models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.contato'),
        ),
        migrations.AddField(
            model_name='cliente',
            name='endereco',
            field=models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, to='sistema_funeraria.endereco'),
        ),
        migrations.AddField(
            model_name='cliente',
            name='user',
            field=models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL),
        ),
    ]
