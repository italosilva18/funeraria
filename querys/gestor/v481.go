package gestor

type Gestor481 struct{}

func (g *Gestor481) TOTAL_VENDIDO_DIA() string {
	return `SELECT sum(IT.valor_real) AS VALOR, sum(DNM.custo) AS CUSTO, sum(IT.valor_real) - sum(DNM.custo) AS LUCRO FROM sm_mv_pdv_it_cup IT, sm_mv_pdv_cb_cup CB, sm_cd_es_produto_dnm DNM WHERE CB.cupom = IT.cupom AND CB.pdv = IT.pdv AND CB.data = IT.data AND CB.empresa = DNM.empresa AND IT.produto = DNM.cod AND CB.data BETWEEN ? AND ? AND IT.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0' AND IT.ic = 0`
}

func (g *Gestor481) TICKET_MEDIO_DIA() string {
	return `SELECT TICKET_MEDIO FROM ( SELECT IT.data, sum(IT.valor_real) / ( SELECT count(DISTINCT(CB.cupom)) FROM sm_mv_pdv_cb_cup CB WHERE CB.data = IT.data AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0' ) AS TICKET_MEDIO FROM sm_mv_pdv_it_cup IT, sm_mv_pdv_cb_cup CP WHERE CP.cupom = IT.cupom AND CP.pdv = IT.pdv AND CP.data = IT.data AND IT.data BETWEEN ? AND ? AND CP.empresa = CAST(1 AS NUMERIC) AND CP.tp = '0' GROUP BY IT.data)`
}

func (g *Gestor481) PRODUTOS_POR_ATENDIMENTO_DIA() string {
	return `SELECT MEDIA_ITENS FROM ( SELECT IT.data, CAST(count(*) AS NUMERIC(15, 2)) / ( SELECT count(DISTINCT(CB.cupom)) FROM sm_mv_pdv_cb_cup CB WHERE CB.data = IT.data AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0') AS MEDIA_ITENS FROM sm_mv_pdv_it_cup IT, sm_mv_pdv_cb_cup CP WHERE IT.cupom = CP.cupom AND IT.pdv = CP.pdv AND IT.data = CP.data AND IT.data BETWEEN ? AND ? AND CP.empresa = CAST(1 AS NUMERIC) AND CP.tp = '0' GROUP BY IT.data)`
}

func (g *Gestor481) TOTAL_CUPONS_VALIDOS_DIA() string {
	return `SELECT count(DISTINCT(CB.cupom)) FROM sm_mv_pdv_cb_cup CB WHERE CB.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0'`
}

func (g *Gestor481) TOTAL_CUPONS_CANCELADOS_DIA() string {
	return `SELECT count(DISTINCT(CB.cupom)) FROM sm_mv_pdv_cb_cup CB WHERE CB.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '1'`
}

func (g *Gestor481) VENDAS_POR_FAIXA_HORARIA() string {
	return `SELECT HORA, VALOR, CUPONS, NUMITENS, CAST( NUMITENS AS NUMERIC(15, 2)) / CUPONS AS MEDIA_ITENS, CAST( VALOR AS NUMERIC(15, 2)) / CUPONS AS MEDIA_VALOR FROM ( SELECT SUBSTRING(CP.hora FROM 1 FOR 2) AS HORA, sum(CP.vlr_liquido) AS VALOR, count (CP.cupom) AS CUPONS, ( SELECT count(IT.produto) FROM sm_mv_pdv_it_cup IT, sm_mv_pdv_cb_cup CB WHERE IT.cupom = CB.cupom AND IT.pdv = CB.pdv AND IT.data = CB.data AND IT.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) AND SUBSTRING(CB.hora FROM 1 FOR 2) = SUBSTRING(CP.hora FROM 1 FOR 2) AND CB.tp = '0') AS NUMITENS FROM sm_mv_pdv_cb_cup CP WHERE CP.data BETWEEN ? AND ? AND CP.empresa = CAST(1 AS NUMERIC) AND CP.tp = '0' GROUP BY SUBSTRING(CP.hora FROM 1 FOR 2) )`
}

func (g *Gestor481) PRODUTOS_MAIS_VENDIDO_DIA() string {
	return `SELECT FIRST 200 EAN, DESCRICAO, TOTAL / QUANTIDADE AS VALOR_UNITARIO, QUANTIDADE, TOTAL, CUSTO, TOTAL - CUSTO AS LUCRO FROM ( SELECT IT.produto AS EAN, PR.dsc AS DESCRICAO, SUM(IT.quantidade) AS QUANTIDADE, SUM(IT.valor_real) TOTAL, SUM(DNM.custo * IT.quantidade) AS CUSTO FROM sm_mv_pdv_it_cup IT, sm_mv_pdv_cb_cup CB, sm_cd_es_produto PR, sm_cd_es_produto_DNM DNM WHERE IT.cupom = CB.cupom AND IT.pdv = CB.pdv AND IT.data = CB.data AND IT.produto = PR.cod AND PR.cod = DNM.cod AND IT.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0' GROUP BY IT.produto, PR.dsc ) ORDER BY TOTAL DESC`
}

func (g *Gestor481) PRODUTOS_ESTATISTICA() string {
	return `SELECT FIRST 5000 EAN, DESCRICAO, TOTAL / QUANTIDADE AS VALOR_UNITARIO, QUANTIDADE, TOTAL, CUSTO, TOTAL - CUSTO AS LUCRO FROM ( SELECT IT.produto AS EAN, PR.dsc AS DESCRICAO, SUM(IT.quantidade) AS QUANTIDADE, SUM(IT.valor_real) TOTAL, SUM(DNM.custo * IT.quantidade) AS CUSTO FROM sm_mv_pdv_it_cup IT, sm_mv_pdv_cb_cup CB, sm_cd_es_produto PR, sm_cd_es_produto_DNM DNM WHERE IT.cupom = CB.cupom AND IT.pdv = CB.pdv AND IT.data = CB.data AND IT.produto = PR.cod AND PR.cod = DNM.cod AND IT.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0' GROUP BY IT.produto, PR.dsc ) ORDER BY TOTAL DESC`
}

func (g *Gestor481) TOTAL_VENDIDO_POR_MODALIDADES() string {
	return `SELECT RD.finalizadora AS CODIGO, FZ.dsc AS DESCRICAO, COUNT(DISTINCT(RD.cupom)) AS QUANTIDADE, SUM(RD.valor - RD.valor_troco) AS TOTAL FROM sm_mv_pdv_cb_cup CB, sm_mv_pdv_rd_cup RD, sm_cd_pdv_finalizadora FZ WHERE CB.cupom = RD.cupom AND CB.pdv = RD.pdv AND CB.data = RD.data AND CB.empresa = FZ.empresa AND RD.finalizadora = FZ.cod AND CB.tp = '0' AND CB.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) GROUP BY RD.finalizadora, FZ.dsc`
}

func (g *Gestor481) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA() string {
	return `SELECT CB.pdv AS PDV, CAST (CAST(CB.data AS DATE) AS VARCHAR(10)) AS DATA, CAST (CAST(CB.hora AS TIME) AS VARCHAR(14)) AS HORA, OP.nome AS OPERADOR, 'SUPERVISOR NÃO IDENTIFICADO' AS SUPERVISOR, CB.vlr_total AS VALOR, CB.cupom AS CUPOM FROM SM_MV_PDV_CB_CUP CB, SM_CD_PDV_OPERADOR OP WHERE CB.operador = OP.cod AND CB.tp = '0' AND CB.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC)`
}

func (g *Gestor481) DESCONTOS_DE_ITENS_DIA() string {
	return `SELECT CAST(CAST(IT.data AS DATE) AS VARCHAR (10)) AS DATA, CAST(CAST(CB.hora AS TIME) AS VARCHAR (14)) AS HORA, OP.nome AS OPERADOR, 'SUPERVISOR DESCONHECIDO' as SUPERVISOR, IT.vlr_ad * -1 AS VALOR, IT.cupom AS CUPOM, IT.item AS ITEM, IT.produto_barras AS EAN, PR.dsc AS DESCRICAO FROM SM_MV_PDV_IT_CUP IT, SM_MV_PDV_CB_CUP CB, SM_CD_ES_PRODUTO PR, SM_CD_PDV_OPERADOR OP WHERE IT.cupom = CB.cupom AND IT.pdv = CB.pdv AND IT.vlr_ad < 0 AND IT.data = CB.data AND IT.produto = PR.cod AND CB.operador = op.cod AND IT.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0' ORDER BY CAST(CB.hora AS TIME)`
}

func (g *Gestor481) VENDAS_POR_SECOES_DIA() string {
	return `SELECT PR.pd_departamento AS CODIGO, DP.dsc AS DESCRICAO, sum(IT.quantidade) AS QUANTIDADE, sum(IT.valor_real) AS VALOR, SUM(IT.quantidade * IT.pr_custo) AS CUSTO, sum(IT.valor_real) - SUM(IT.quantidade * IT.pr_custo) AS LUCRO FROM SM_MV_PDV_IT_CUP IT, SM_MV_PDV_CB_CUP CB, SM_CD_ES_PRODUTO PR, SM_CD_ES_DEPARTAMENTO DP WHERE IT.cupom = CB.cupom AND IT.pdv = CB.pdv AND IT.data = CB.data AND IT.produto = PR.cod AND PR.pd_departamento = DP.cod AND IT.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0' AND IT.ic = 0 GROUP BY PR.pd_departamento, DP.dsc`
}

func (g *Gestor481) TOTAL_CUPONS_POR_OPERADOR_DIA() string {
	return `SELECT COUNT(CB.cupom) AS QUANTIDADE, OP.nome AS OPERADOR, SUM(CB.vlr_total) AS VALOR FROM SM_MV_PDV_CB_CUP CB, SM_CD_PDV_OPERADOR OP WHERE CB.operador = OP.cod AND CB.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) AND CB.tp = '0' GROUP BY OP.nome`
}

func (g *Gestor481) VENDAS_POR_VENDEDOR_DIA() string {
	return ""
}

func (g *Gestor481) MODALIDADES_PAGAMENTO_OPERADOR_DIA() string {
	return `SELECT OP.nome AS OPERADOR, RD.finalizadora AS CODIGO, FZ.dsc AS DESCRICAO, COUNT(DISTINCT(RD.cupom)) AS QUANTIDADE, SUM(RD.valor-RD.valor_troco) AS TOTAL FROM sm_mv_pdv_cb_cup CB, sm_mv_pdv_rd_cup RD, sm_cd_pdv_finalizadora FZ, sm_cd_pdv_operador OP WHERE CB.cupom = RD.cupom AND CB.PDV = RD.pdv AND CB.data = RD.data AND RD.finalizadora = FZ.cod AND CB.operador = OP.cod AND CB.tp = '0' AND CB.data BETWEEN ? AND ? AND CB.empresa = CAST(1 AS NUMERIC) GROUP BY RD.finalizadora, FZ.dsc, OP.nome`
}
