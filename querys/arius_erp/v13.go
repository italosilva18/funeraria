package arius_erp

type AriusErp13 struct{}

// TOTAL_VENDIDO_DIA - Alterada 01-02-2019
func (g *AriusErp13) TOTAL_VENDIDO_DIA() string {
	return "select sum(iv.valor - iv.desconto + iv.acrescimo) as valor, sum(pl.custo * iv.qtde) as custo, (sum(iv.valor - iv.desconto + iv.acrescimo) - sum(pl.custo * iv.qtde)) as lucro from vendas v inner join itens_venda iv on v.id = iv.venda inner join produtos_loja pl on v.empresa = pl.politica and iv.produto = pl.id where trunc(v.data_hora) = TO_DATE(:DATA, 'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and nvl(v.devolucao,'F') = 'F' and iv.valor > 0 and nvl(v.cancelado,'F') = 'F'"
}

// TICKET_MEDIO_DIA - Alterada 01-02-2019
func (g *AriusErp13) TICKET_MEDIO_DIA() string {
	return "select round(sum(v.valor) / count(v.id),4) as ticket_medio from vendas v where trunc(v.data_hora) = TO_DATE(:DATA, 'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and (nvl(v.devolucao,'F') = 'F' or v.valor > 0) and nvl(v.cancelado,'F') = 'F'"
}

// PRODUTOS_POR_ATENDIMENTO_DIA - Alterada 01-02-2019
func (g *AriusErp13) PRODUTOS_POR_ATENDIMENTO_DIA() string {
	return "select round(count(iv.produto) / count(distinct iv.venda),4) numero_itens from vendas v inner join itens_venda iv on v.id = iv.venda where trunc(v.data_hora) = TO_DATE(:DATA,'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and (nvl(v.devolucao,'F') = 'F' or v.valor > 0) and nvl(v.cancelado, 'F') = 'F' and iv.valor > 0"
}

// TOTAL_CUPONS_VALIDOS_DIA - Alterada 01-02-2019
func (g *AriusErp13) TOTAL_CUPONS_VALIDOS_DIA() string {
	return "select count(v.id) as cupons from vendas v where trunc(v.data_hora) = TO_DATE(:DATA, 'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and nvl(v.devolucao,'F') = 'F' and nvl(v.cancelado,'F') = 'F' AND v.valor > 0"
}

// TOTAL_CUPONS_CANCELADO_DIAS - Alterada 01-02-2019
func (g *AriusErp13) TOTAL_CUPONS_CANCELADOS_DIA() string {
	return "select count(v.id) as cupons from vendas v where trunc(v.data_hora) = TO_DATE(:DATA,'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and (nvl(v.devolucao,'F') = 'T' or v.valor < 0 or nvl(v.cancelado,'F') = 'T' ) -- Devolucoes do ERP / PDV"
}

// VENDAS_POR_FAIXA_HORARIA - Alterada 01-02-2019
func (g *AriusErp13) VENDAS_POR_FAIXA_HORARIA() string {
	return "select TO_CHAR(trunc(v.data_hora, 'HH24'),'HH24')as hora, sum(iv.valor - iv.desconto + iv.acrescimo) valor, count(distinct iv.venda) as numero_cupom, count(iv.venda) as numero_itens, round(count(iv.id) / count(distinct v.id),4) as media_itens, round(sum(iv.valor - iv.desconto + iv.acrescimo) / count(distinct v.id),4) as media_valor from vendas v inner join itens_venda iv on v.id = iv.venda where v.empresa = :ID and trunc(v.data_hora) = to_date(:DATA, 'dd/mm/rrrr') and nvl(v.devolucao,'F') = 'F' and nvl(v.cancelado,'F') = 'F' and v.valor > 0 group by trunc(v.data_hora, 'HH24')"
}

// PRODUTOS_MAIS_VENDIDOS_DIA - Alterada 01-02-2019
func (g *AriusErp13) PRODUTOS_MAIS_VENDIDO_DIA() string {
	return "select * from (select nvl(iv.ean, iv.produto) as ean, p.descritivo as descricao, round((fn_zero_to_one(sum(iv.valor - iv.desconto + iv.acrescimo)) / fn_zero_to_one(sum(iv.qtde))),4) as valor_unitario, sum(iv.qtde) as quantidade, sum(iv.valor - iv.desconto + iv.acrescimo) as total, sum(iv.qtde * pl.custo) as custo, (sum(iv.valor - iv.desconto + iv.acrescimo) - sum(iv.qtde * pl.custo)) as lucro from vendas v inner join itens_venda iv on v.id = iv.venda inner join produtos p on iv.produto = p.id inner join produtos_loja pl on v.empresa = pl.politica and iv.produto = pl.id where v.empresa = TO_NUMBER(:ID) and trunc(v.data_hora) = TO_DATE(:DATA,'DD/MM/YYYY') and nvl(v.devolucao,'F') = 'F' and v.valor > 0 and nvl(v.cancelado,'F') = 'F' group by nvl(iv.ean, iv.produto), p.descritivo order by 5 desc) where rownum <= 200"
}

// PRODUTOS_ESTATISTICA - Alterada 01-02-2019
func (g *AriusErp13) PRODUTOS_ESTATISTICA() string {
	return "select * from (select nvl(iv.ean, iv.produto) as ean, p.descritivo as descricao, round((fn_zero_to_one(sum(iv.valor)) / fn_zero_to_one(sum(iv.qtde))),4) as valor_unitario, sum(iv.qtde) as quantidade, sum(iv.valor) as total, sum(iv.qtde * pl.custo) as custo, (sum(iv.valor) - sum(iv.qtde * pl.custo)) as lucro from vendas v inner join itens_venda iv on v.id = iv.venda inner join produtos p on iv.produto = p.id inner join produtos_loja pl on v.empresa = pl.politica and iv.produto = pl.id where v.empresa = TO_NUMBER(:ID) and trunc(v.data_hora) = TO_DATE(:DATA, 'DD/MM/YYYY') AND V.EMPRESA = TO_NUMBER(:ID) and nvl(v.devolucao,'F') = 'F' and v.valor > 0 and nvl(v.cancelado,'F') = 'F' group by nvl(iv.ean, iv.produto), p.descritivo order by 5 desc) where rownum <= 8000"
}

// TOTAL_VENDIDO_POR_MODALIDADES - Alterada 01-02-2019
func (g *AriusErp13) TOTAL_VENDIDO_POR_MODALIDADES() string {
	return "select nvl(f.forma_pagto,99) as codigo, nvl(fp.descritivo,'DEVOLUCAO') as descricao, count(f.id) as quantidade, sum(f.valor) as valor from vendas v left join finalizacoes f on v.id = f.venda left join formas_pagto fp on f.forma_pagto = fp.id where v.empresa = TO_NUMBER(:ID) and trunc(v.data_hora) = TO_DATE(:DATA,'DD/MM/YYYY') and nvl(v.devolucao,'F') = 'F' and nvl(v.cancelado,'F') = 'F' and v.valor > 0 group by f.forma_pagto, fp.descritivo order by valor desc"
}

// ESTORNO_DE_CUPONS_POR_OPERADOR_DIA - Alterada 01-02-2019
func (g *AriusErp13) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA() string {
	return "select v.pdv, to_char(v.data_hora, 'dd/mm/yyyy') as data, to_char(v.data_hora, 'HH24:MI') as hora, f.nome as operador, nvl(f.nome, null) as operador, abs(v.valor) as valor, TO_CHAR(nvl(v.nf, v.documento)) as cupom from vendas v left join bas_t_func (g *AriusErp13)ionarios f on v.operador = f.id_func (g *AriusErp13)ionario where trunc(v.data_hora) = TO_DATE(:DATA,'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and (v.valor < 0 or v.cancelado = 'T' or v.devolucao = 'T')"
}

// DESCONTOS_DE_ITENS_DIA - Alterada 01-02-2019
func (g *AriusErp13) DESCONTOS_DE_ITENS_DIA() string {
	return "select to_char(v.data_hora, 'dd/mm/yyyy') as data, to_char(v.data_hora, 'HH24:MI') as hora, nvl(opr.nome,null) as operador, nvl(opr.nome,null) as supervisor, iv.desconto as valor, nvl(v.documento, v.nf) as cupom, rownum  as item, nvl(iv.ean, iv.produto) as ean, p.descritivo as descricao from vendas v inner join itens_venda iv on v.id = iv.venda inner join produtos p on iv.produto = p.id left join produtos_ean pe on p.id = pe.produto and rownum = 1 and pe.qtdee = 1 left join bas_t_func (g *AriusErp13)ionarios opr on v.operador = opr.id_func (g *AriusErp13)ionario where v.empresa = TO_NUMBER(&ID) and trunc(v.data_hora) = to_date(&DATA ,'dd/mm/rrrr') and nvl(v.devolucao,'F') = 'F' and nvl(v.cancelado,'F') = 'F' and v.valor > 0 and iv.desconto > 0"
}

// VENDAS_POR_SECOES_DIA - Alterada 01-02-2019
func (g *AriusErp13) VENDAS_POR_SECOES_DIA() string {
	return "select d.depto as codigo, d.descritivo as descricao, count(distinct iv.produto) as qtde, sum(iv.valor - iv.desconto + iv.acrescimo) as valor, sum(iv.qtde * pl.custo) as cmv, (sum(iv.valor - iv.desconto + iv.acrescimo) - sum(iv.qtde * pl.custo)) as lucro from vendas v inner join itens_venda iv on v.id = iv.venda inner join produtos p on iv.produto = p.id inner join deptos d on p.depto = d.depto and d.secao = 0 inner join produtos_loja pl on v.empresa = pl.politica and iv.produto = pl.id where trunc(v.data_hora) = TO_DATE(:DATA, 'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and nvl(v.devolucao,'F') = 'F' and nvl(v.cancelado,'F') = 'F' and iv.valor > 0 group by d.depto, d.descritivo order by valor desc"
}

// TOTAL_CUPONS_POR_OPERADOR_DIA - Alterada 01-02-2019
func (g *AriusErp13) TOTAL_CUPONS_POR_OPERADOR_DIA() string {
	return "select count(v.id) as qtde_cupons, nvl(f.nome, null) as operador, sum(v.valor) as valor from vendas v left join bas_t_func (g *AriusErp13)ionarios f on v.operador = f.id_func (g *AriusErp13)ionario where trunc(v.data_hora) = TO_DATE(:DATA,'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and nvl(v.devolucao,'F') = 'F' and nvl(v.cancelado,'F') = 'F' and v.valor > 0 group by nvl(f.nome, null) order by qtde_cupons desc"
}

// VENDAS_POR_VENDEDOR_DIA - Alterada 01-02-2019
func (g *AriusErp13) VENDAS_POR_VENDEDOR_DIA() string {
	return "select nvl(f.descritivo, 'null') as vendedor, sum(iv.valor - iv.desconto + iv.acrescimo) as valor, count(distinct v.id) as cupons, count(iv.produto) as item, round(sum(count(iv.produto)) over (order by null rows between unbounded preceding and unbounded following) / count(iv.produto),2) as media_item, round(sum(iv.valor - iv.desconto + iv.acrescimo) / count(iv.produto),4) as media_valor from vendas v inner join itens_venda iv on v.id = iv.venda left join vendedores f on iv.vendedor = f.id where trunc(v.data_hora) = TO_DATE(:DATA, 'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and nvl(v.devolucao,'F') = 'F' and nvl(v.cancelado,'F') = 'F' and v.valor > 0 group by trunc(v.data_hora), nvl(f.descritivo, 'null'), iv.vendedor"
}

// MODALIDADES_PAGAMENTO_OPERADOR_DIA - Alterada 01-02-2019
func (g *AriusErp13) MODALIDADES_PAGAMENTO_OPERADOR_DIA() string {
	return "select TO_CHAR(trunc(v.data_hora),'dd/mm/yyyy') as data, nvl(opr.nome,null) as operador, f.forma_pagto as cod_pagto, fp.descritivo as descricao, count(distinct v.id) as qtde, sum(iv.valor - iv.desconto + iv.acrescimo) as valor from vendas v inner join itens_venda iv on v.id = iv.venda inner join finalizacoes f on v.id = f.venda inner join formas_pagto fp on f.forma_pagto = fp.id left join bas_t_func (g *AriusErp13)ionarios opr on v.operador = opr.id_func (g *AriusErp13)ionario where trunc(v.data_hora) = TO_DATE(:DATA,'DD/MM/YYYY') and v.empresa = TO_NUMBER(:ID) and nvl(v.devolucao,'F') = 'F' and nvl(v.cancelado,'F') = 'F' and v.valor > 0 group by trunc(v.data_hora), nvl(opr.nome,null), f.forma_pagto, fp.descritivo"
}
