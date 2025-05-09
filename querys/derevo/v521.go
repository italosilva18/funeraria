package derevo

type Derevo521 struct{}

// 1 - Total vendido (diário)
func (d *Derevo521) TOTAL_VENDIDO_DIA() string {
	return `
SELECT
    COALESCE(C1.pvalor, 0) AS valor,
    COALESCE(C1.pcusto, 0) AS custo,
    COALESCE(C1.plucro_bruto, 0) AS lucro
FROM SP_SM_PS_BI_PDV_CUPOM_ITEM_C1(?, ?, ?) C1;
`
}

// 2 - Ticket médio (diário)
func (d *Derevo521) TICKET_MEDIO_DIA() string {
	return `
SELECT
    COALESCE(C3.pticket_medio, 0) AS ticket_medio
FROM SP_SM_PS_BI_PDV_CUPOM_C3(?, ?, ?) C3;
`
}

// 3 - Média de produtos por atendimento (diário)
func (d *Derevo521) PRODUTOS_POR_ATENDIMENTO_DIA() string {
	return `
SELECT
    COALESCE(C3.pmedia_item, 0) AS media_itens
FROM SP_SM_PS_BI_PDV_CUPOM_ITEM_C3(?, ?, ?) C3;
`
}

// 4 - Cupons válidos (diário)
func (d *Derevo521) TOTAL_CUPONS_VALIDOS_DIA() string {
	return `
SELECT
    COALESCE(C2.pqnt_cupom, 0) AS cupons
FROM SP_SM_PS_BI_PDV_CUPOM_C2(?, ?, ?) C2;
`
}

// 5 - Cupons cancelados (diário)
func (d *Derevo521) TOTAL_CUPONS_CANCELADOS_DIA() string {
	return `
SELECT
    COALESCE(C2.pqnt_cupom_cancelado, 0) AS cancelados
FROM SP_SM_PS_BI_PDV_CUPOM_C2(?, ?, ?) C2;
`
}

// 6 - Vendas por faixa horária (diário)
func (d *Derevo521) VENDAS_POR_FAIXA_HORARIA() string {
	return `
SELECT
    C4.phora AS hora,
    COALESCE(C4.pvalor, 0) AS valor,
    COALESCE(C4.pqnt_cupom, 0) AS num_cupons,
    COALESCE(C4.pqnt_item, 0) AS num_itens,
    COALESCE(C4.pmedia_itens, 0) AS media_itens,
    COALESCE(C4.pmedia_valor, 0) AS media_valor
FROM SP_SM_PS_BI_PDV_CUPOM_ITEM_C4(?, ?, ?) C4;
`
}

// 7 - Produtos mais vendidos (100 primeiros)
func (d *Derevo521) PRODUTOS_MAIS_VENDIDO_DIA() string {
	return `
SELECT
    IIF(C5.pproduto_b_tipo = 0 AND C5.pproduto_b > 0, C5.pproduto_b, C5.pproduto) AS ean,
    C5.pproduto_d AS descricao,
    C5.pvalor_unitario AS valor_unitario,
    C5.pqnt AS quantidade,
    C5.ptotal AS total,
    C5.pcusto AS custo,
    C5.plucro_bruto AS lucro
FROM SP_SM_PS_BI_PDV_CUPOM_ITEM_C5(?, ?, ?) C5;
`
}

// 8 - Produtos estatística (5000 produtos)
func (d *Derevo521) PRODUTOS_ESTATISTICA() string {
	return `
SELECT
    IIF(C6.pproduto_b_tipo = 0 AND C6.pproduto_b > 0, C6.pproduto_b, C6.pproduto) AS ean,
    C6.pproduto_d AS descricao,
    C6.pvalor_unitario AS valor_unitario,
    C6.pqnt AS quantidade,
    C6.ptotal AS total,
    C6.pcusto AS custo,
    C6.plucro_bruto AS lucro
FROM SP_SM_PS_BI_PDV_CUPOM_ITEM_C6(?, ?, ?) C6;
`
}

// 9 - Total vendido por modalidade
func (d *Derevo521) TOTAL_VENDIDO_POR_MODALIDADES() string {
	return `
SELECT
    C1.pfinalizadora AS codigo,
    C1.pfinalizadora_d AS descricao,
    C1.pqnt AS quantidade,
    C1.ptotal AS total
FROM SP_SM_PS_BI_PDV_CUPOM_FNZ_C1(?, ?, ?) C1;
`
}

// 10 - Estorno de cupons
func (d *Derevo521) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA() string {
	return `
SELECT
    C4.ppdv AS pdv,
    C4.pdata AS data,
    C4.phora AS hora,
    C4.poperador_d AS operador,
    C4.poperador_aut_d AS supervisor,
    C4.pvalor_liquido,
    C4.pcupom
FROM SP_SM_PS_BI_PDV_CUPOM_C4(?, ?, ?) C4;
`
}

// 11 - Descontos de itens
func (d *Derevo521) DESCONTOS_DE_ITENS_DIA() string {
	return `
SELECT
    C7.pdata AS data,
    C7.phora AS hora,
    C7.poperador_d AS operador,
    C7.poperador_aut_d AS supervisor,
    C7.pvalor_desconto AS valor,
    C7.pcupom AS cupom,
    C7.pitem AS item,
    IIF(C7.pproduto_b_tipo = 0 AND C7.pproduto_b > 0, C7.pproduto_b, C7.pproduto) AS ean,
    C7.pproduto_d AS descricao
FROM SP_SM_PS_BI_PDV_CUPOM_ITEM_C7(?, ?, ?) C7;
`
}

// 12 - Vendas por seções
func (d *Derevo521) VENDAS_POR_SECOES_DIA() string {
	return `
SELECT
    C9.pdepartamento AS codigo,
    C9.pdepartamento_d AS descricao,
    C9.pqnt AS quantidade,
    C9.pvalor AS valor,
    C9.pcusto AS custo,
    C9.plucro_bruto AS lucro
FROM SP_SM_PS_BI_PDV_CUPOM_ITEM_C9(?, ?, ?) C9;
`
}

// 13 - Total de cupons por operador
func (d *Derevo521) TOTAL_CUPONS_POR_OPERADOR_DIA() string {
	return `
SELECT
    C5.pqnt AS qtd_cupons,
    C5.poperador_d AS operador,
    C5.pvalor AS valor
FROM SP_SM_PS_BI_PDV_CUPOM_C5(?, ?, ?) C5;
`
}

// 14 - Vendas por vendedor (caso não tenha, retorne string vazia)
func (d *Derevo521) VENDAS_POR_VENDEDOR_DIA() string {
	return ""
}

// 15 - Modalidades por operador
func (d *Derevo521) MODALIDADES_PAGAMENTO_OPERADOR_DIA() string {
	return `
SELECT
    C2.pdata AS data,
    C2.poperador_d AS operador,
    C2.pfinalizadora_d AS cod_pagto,
    C2.pqnt AS qtd,
    C2.ptotal AS valor
FROM SP_SM_PS_BI_PDV_CUPOM_FNZ_C2(?, ?, ?) C2;
`
}
