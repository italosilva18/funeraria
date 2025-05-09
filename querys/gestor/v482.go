package gestor

type Gestor482 struct{}

// TOTAL_VENDIDO_DIA - Retorna três FLOAT: Valor, Custo e Lucro.
func (g *Gestor482) TOTAL_VENDIDO_DIA() string {
	return `
SELECT 
    CAST(SUM(IT.valor_real) AS FLOAT) AS VALOR,
    CAST(SUM(DNM.custo) AS FLOAT) AS CUSTO,
    CAST(SUM(IT.valor_real) - SUM(DNM.custo) AS FLOAT) AS LUCRO
FROM sm_mv_pdv_it_cup IT,
     sm_mv_pdv_cb_cup CB,
     sm_cd_es_produto_dnm DNM
WHERE CB.cupom = IT.cupom 
  AND CB.pdv = IT.pdv 
  AND CB.data = IT.data 
  AND CB.empresa = DNM.empresa 
  AND IT.produto = DNM.cod 
  AND CB.data BETWEEN ? AND ?
  AND IT.empresa = ? 
  AND CB.tp = '0' 
  AND IT.ic = 0`
}

// TICKET_MEDIO_DIA - Retorna um FLOAT.
func (g *Gestor482) TICKET_MEDIO_DIA() string {
	return `
SELECT 
    CAST(TICKET_MEDIO AS FLOAT) AS TICKET_MEDIO
FROM (
    SELECT IT.data,
           SUM(IT.valor_real) / 
           (SELECT COUNT(DISTINCT CB.cupom)
            FROM sm_mv_pdv_cb_cup CB
            WHERE CB.data = IT.data 
              AND CB.empresa = ? 
              AND CB.tp = '0'
           ) AS TICKET_MEDIO
    FROM sm_mv_pdv_it_cup IT,
         sm_mv_pdv_cb_cup CP
    WHERE CP.cupom = IT.cupom 
      AND CP.pdv = IT.pdv 
      AND CP.data = IT.data 
      AND IT.data BETWEEN ? AND ?
      AND CP.empresa = ? 
      AND CP.tp = '0'
    GROUP BY IT.data
) sub`
}

// PRODUTOS_POR_ATENDIMENTO_DIA - Retorna um FLOAT.
func (g *Gestor482) PRODUTOS_POR_ATENDIMENTO_DIA() string {
	return `
SELECT 
    CAST(MEDIA_ITENS AS FLOAT) AS MEDIA_ITENS
FROM (
    SELECT IT.data,
           CAST(COUNT(*) AS NUMERIC(15,2)) / 
           (SELECT COUNT(DISTINCT CB.cupom)
            FROM sm_mv_pdv_cb_cup CB
            WHERE CB.data = IT.data 
              AND CB.empresa = ? 
              AND CB.tp = '0'
           ) AS MEDIA_ITENS
    FROM sm_mv_pdv_it_cup IT,
         sm_mv_pdv_cb_cup CP
    WHERE IT.cupom = CP.cupom 
      AND IT.pdv = CP.pdv 
      AND IT.data = CP.data 
      AND IT.data BETWEEN ? AND ?
      AND CP.empresa = ? 
      AND CP.tp = '0'
    GROUP BY IT.data
) sub`
}

// TOTAL_CUPONS_VALIDOS_DIA - Retorna um INTEGER.
func (g *Gestor482) TOTAL_CUPONS_VALIDOS_DIA() string {
	return `
SELECT 
    CAST(COUNT(DISTINCT CB.cupom) AS INTEGER) AS TOTAL_CUPONS
FROM sm_mv_pdv_cb_cup CB
WHERE CB.data BETWEEN ? AND ?
  AND CB.empresa = ? 
  AND CB.tp = '0'`
}

// TOTAL_CUPONS_CANCELADOS_DIA - Retorna um INTEGER.
func (g *Gestor482) TOTAL_CUPONS_CANCELADOS_DIA() string {
	return `
SELECT 
    CAST(COUNT(DISTINCT CB.cupom) AS INTEGER) AS TOTAL_CUPONS_CANCELADOS
FROM sm_mv_pdv_cb_cup CB
WHERE CB.data BETWEEN ? AND ?
  AND CB.empresa = ? 
  AND CB.tp = '1'`
}

// VENDAS_POR_FAIXA_HORARIA - Retorna: HORA (INTEGER), VALOR (FLOAT), CUPONS (INTEGER), NUMITENS (INTEGER),
// MEDIA_ITENS (FLOAT) e MEDIA_VALOR (FLOAT).
func (g *Gestor482) VENDAS_POR_FAIXA_HORARIA() string {
	return `
SELECT 
    CAST(HORA AS INTEGER) AS HORA,
    CAST(VALOR AS FLOAT) AS VALOR,
    CAST(CUPONS AS INTEGER) AS CUPONS,
    CAST(NUMITENS AS INTEGER) AS NUMITENS,
    CAST(CAST(NUMITENS AS NUMERIC(15,2)) / CUPONS AS FLOAT) AS MEDIA_ITENS,
    CAST(CAST(VALOR AS NUMERIC(15,2)) / CUPONS AS FLOAT) AS MEDIA_VALOR
FROM (
    SELECT 
        SUBSTRING(CP.hora FROM 1 FOR 2) AS HORA,
        SUM(CP.vlr_liquido) AS VALOR,
        COUNT(CP.cupom) AS CUPONS,
        (SELECT COUNT(IT.produto)
         FROM sm_mv_pdv_it_cup IT,
              sm_mv_pdv_cb_cup CB
         WHERE IT.cupom = CB.cupom 
           AND IT.pdv = CB.pdv 
           AND IT.data = CB.data
           AND IT.data BETWEEN ? AND ?
           AND CB.empresa = ?
           AND SUBSTRING(CB.hora FROM 1 FOR 2) = SUBSTRING(CP.hora FROM 1 FOR 2)
           AND CB.tp = '0'
        ) AS NUMITENS
    FROM sm_mv_pdv_cb_cup CP
    WHERE CP.data BETWEEN ? AND ?
      AND CP.empresa = ? 
      AND CP.tp = '0'
    GROUP BY SUBSTRING(CP.hora FROM 1 FOR 2)
) sub`
}

// PRODUTOS_MAIS_VENDIDO_DIA - Retorna: EAN, DESCRICAO, VALOR_UNITARIO (FLOAT), QUANTIDADE (INTEGER),
// TOTAL (FLOAT), CUSTO (FLOAT) e LUCRO (FLOAT).
func (g *Gestor482) PRODUTOS_MAIS_VENDIDO_DIA() string {
	return `
SELECT FIRST 200 
    EAN,
    DESCRICAO,
    CAST(TOTAL / QUANTIDADE AS FLOAT) AS VALOR_UNITARIO,
    CAST(QUANTIDADE AS INTEGER) AS QUANTIDADE,
    CAST(TOTAL AS FLOAT) AS TOTAL,
    CAST(CUSTO AS FLOAT) AS CUSTO,
    CAST(TOTAL - CUSTO AS FLOAT) AS LUCRO
FROM (
    SELECT 
        IT.produto AS EAN,
        PR.dsc AS DESCRICAO,
        SUM(IT.quantidade) AS QUANTIDADE,
        SUM(IT.valor_real) AS TOTAL,
        SUM(DNM.custo * IT.quantidade) AS CUSTO
    FROM sm_mv_pdv_it_cup IT,
         sm_mv_pdv_cb_cup CB,
         sm_cd_es_produto PR,
         sm_cd_es_produto_DNM DNM
    WHERE IT.cupom = CB.cupom 
      AND IT.pdv = CB.pdv 
      AND IT.data = CB.data 
      AND IT.produto = PR.cod 
      AND PR.cod = DNM.cod 
      AND IT.data BETWEEN ? AND ?
      AND CB.empresa = ? 
      AND CB.tp = '0'
    GROUP BY IT.produto, PR.dsc
) sub
ORDER BY TOTAL DESC`
}

// PRODUTOS_ESTATISTICA - Mesmos campos e conversões da query anterior para os 5000 itens.
func (g *Gestor482) PRODUTOS_ESTATISTICA() string {
	return `
SELECT FIRST 5000 
    EAN,
    DESCRICAO,
    CAST(TOTAL / QUANTIDADE AS FLOAT) AS VALOR_UNITARIO,
    CAST(QUANTIDADE AS INTEGER) AS QUANTIDADE,
    CAST(TOTAL AS FLOAT) AS TOTAL,
    CAST(CUSTO AS FLOAT) AS CUSTO,
    CAST(TOTAL - CUSTO AS FLOAT) AS LUCRO
FROM (
    SELECT 
        IT.produto AS EAN,
        PR.dsc AS DESCRICAO,
        SUM(IT.quantidade) AS QUANTIDADE,
        SUM(IT.valor_real) AS TOTAL,
        SUM(DNM.custo * IT.quantidade) AS CUSTO
    FROM sm_mv_pdv_it_cup IT,
         sm_mv_pdv_cb_cup CB,
         sm_cd_es_produto PR,
         sm_cd_es_produto_DNM DNM
    WHERE IT.cupom = CB.cupom 
      AND IT.pdv = CB.pdv 
      AND IT.data = CB.data 
      AND IT.produto = PR.cod 
      AND PR.cod = DNM.cod 
      AND IT.data BETWEEN ? AND ?
      AND CB.empresa = ? 
      AND CB.tp = '0'
    GROUP BY IT.produto, PR.dsc
) sub
ORDER BY TOTAL DESC`
}

// TOTAL_VENDIDO_POR_MODALIDADES - Retorna: CODIGO (INTEGER), DESCRICAO (STRING),
// QUANTIDADE (INTEGER) e TOTAL (FLOAT).
func (g *Gestor482) TOTAL_VENDIDO_POR_MODALIDADES() string {
	return `
SELECT 
    CAST(RD.finalizadora AS INTEGER) AS CODIGO,
    FZ.dsc AS DESCRICAO,
    CAST(COUNT(DISTINCT RD.cupom) AS INTEGER) AS QUANTIDADE,
    CAST(SUM(RD.valor - RD.valor_troco) AS FLOAT) AS TOTAL
FROM sm_mv_pdv_cb_cup CB,
     sm_mv_pdv_rd_cup RD,
     sm_cd_pdv_finalizadora FZ
WHERE CB.cupom = RD.cupom 
  AND CB.pdv = RD.pdv 
  AND CB.data = RD.data 
  AND CB.empresa = FZ.empresa 
  AND RD.finalizadora = FZ.cod 
  AND CB.tp = '0'
  AND CB.data BETWEEN ? AND ?
  AND CB.empresa = ?
GROUP BY RD.finalizadora, FZ.dsc`
}

// ESTORNO_DE_CUPONS_POR_OPERADOR_DIA - Retorna: PDV (STRING), DATA (STRING),
// HORA (STRING), OPERADOR (STRING), SUPERVISOR (STRING), VALOR (FLOAT) e CUPOM (STRING).
func (g *Gestor482) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA() string {
	return `
SELECT 
    CB.pdv AS PDV,
    CAST(CAST(CB.data AS DATE) AS VARCHAR(10)) AS DATA,
    CAST(CAST(CB.hora AS TIME) AS VARCHAR(14)) AS HORA,
    OP.nome AS OPERADOR,
    'SUPERVISOR NÃO IDENTIFICADO' AS SUPERVISOR,
    CAST(CB.vlr_total AS FLOAT) AS VALOR,
    CB.cupom AS CUPOM
FROM SM_MV_PDV_CB_CUP CB,
     SM_CD_PDV_OPERADOR OP
WHERE CB.operador = OP.cod 
  AND CB.tp = '0'
  AND CB.data BETWEEN ? AND ?
  AND CB.empresa = ?`
}

// DESCONTOS_DE_ITENS_DIA - Retorna: DATA (STRING), HORA (STRING), OPERADOR (STRING),
// SUPERVISOR (STRING), VALOR (FLOAT), CUPOM (STRING), ITEM (INTEGER), EAN (STRING) e DESCRICAO (STRING).
func (g *Gestor482) DESCONTOS_DE_ITENS_DIA() string {
	return `
SELECT 
    CAST(CAST(IT.data AS DATE) AS VARCHAR(10)) AS DATA,
    CAST(CAST(CB.hora AS TIME) AS VARCHAR(14)) AS HORA,
    OP.nome AS OPERADOR,
    'SUPERVISOR DESCONHECIDO' AS SUPERVISOR,
    CAST(IT.VLR_DESCONTO * -1 AS FLOAT) AS VALOR,
    IT.CUPOM AS CUPOM,
    CAST(IT.item AS INTEGER) AS ITEM,
    IT.produto_barras AS EAN,
    PR.dsc AS DESCRICAO
FROM SM_MV_PDV_IT_CUP IT,
     SM_MV_PDV_CB_CUP CB,
     SM_CD_ES_PRODUTO PR,
     SM_CD_PDV_OPERADOR OP
WHERE IT.CUPOM = CB.cupom 
  AND IT.pdv = CB.pdv 
  AND IT.VLR_DESCONTO > 0 
  AND IT.data = CB.data 
  AND IT.produto = PR.cod 
  AND CB.operador = OP.cod 
  AND IT.data BETWEEN ? AND ?
  AND CB.empresa = ? 
  AND CB.tp = '0'
ORDER BY CAST(CB.hora AS TIME)`
}

// VENDAS_POR_SECOES_DIA - Retorna: CODIGO (INTEGER), DESCRICAO (STRING),
// QUANTIDADE (INTEGER), VALOR (FLOAT), CUSTO (FLOAT) e LUCRO (FLOAT).
func (g *Gestor482) VENDAS_POR_SECOES_DIA() string {
	return `
SELECT 
    CAST(PR.pd_departamento AS INTEGER) AS CODIGO,
    DP.dsc AS DESCRICAO,
    CAST(SUM(IT.quantidade) AS INTEGER) AS QUANTIDADE,
    CAST(SUM(IT.valor_real) AS FLOAT) AS VALOR,
    CAST(SUM(IT.quantidade * IT.pr_custo) AS FLOAT) AS CUSTO,
    CAST(SUM(IT.valor_real) - SUM(IT.quantidade * IT.pr_custo) AS FLOAT) AS LUCRO
FROM SM_MV_PDV_IT_CUP IT,
     SM_MV_PDV_CB_CUP CB,
     SM_CD_ES_PRODUTO PR,
     SM_CD_ES_DEPARTAMENTO DP
WHERE IT.cupom = CB.cupom 
  AND IT.pdv = CB.pdv 
  AND IT.data = CB.data 
  AND IT.produto = PR.cod 
  AND PR.pd_departamento = DP.cod 
  AND IT.data BETWEEN ? AND ?
  AND CB.empresa = ? 
  AND CB.tp = '0'
  AND IT.ic = 0
GROUP BY PR.pd_departamento, DP.dsc`
}

// TOTAL_CUPONS_POR_OPERADOR_DIA - Retorna: QUANTIDADE (INTEGER), OPERADOR (STRING) e VALOR (FLOAT).
func (g *Gestor482) TOTAL_CUPONS_POR_OPERADOR_DIA() string {
	return `
SELECT 
    CAST(COUNT(CB.cupom) AS INTEGER) AS QUANTIDADE,
    OP.nome AS OPERADOR,
    CAST(SUM(CB.vlr_total) AS FLOAT) AS VALOR
FROM SM_MV_PDV_CB_CUP CB,
     SM_CD_PDV_OPERADOR OP
WHERE CB.operador = OP.cod 
  AND CB.data BETWEEN ? AND ?
  AND CB.empresa = ? 
  AND CB.tp = '0'
GROUP BY OP.nome`
}

// VENDAS_POR_VENDEDOR_DIA - Deixa-se vazio conforme modelo.
func (g *Gestor482) VENDAS_POR_VENDEDOR_DIA() string {
	return ""
}

// MODALIDADES_PAGAMENTO_OPERADOR_DIA - Retorna: OPERADOR (STRING), CODIGO (INTEGER),
// DESCRICAO (STRING), QUANTIDADE (INTEGER) e TOTAL (FLOAT).
func (g *Gestor482) MODALIDADES_PAGAMENTO_OPERADOR_DIA() string {
	return `
SELECT 
    OP.nome AS OPERADOR,
    CAST(RD.finalizadora AS INTEGER) AS CODIGO,
    FZ.dsc AS DESCRICAO,
    CAST(COUNT(DISTINCT RD.cupom) AS INTEGER) AS QUANTIDADE,
    CAST(SUM(RD.valor - RD.valor_troco) AS FLOAT) AS TOTAL
FROM sm_mv_pdv_cb_cup CB,
     sm_mv_pdv_rd_cup RD,
     sm_cd_pdv_finalizadora FZ,
     SM_CD_PDV_OPERADOR OP
WHERE CB.cupom = RD.cupom 
  AND CB.PDV = RD.pdv 
  AND CB.data = RD.data 
  AND RD.finalizadora = FZ.cod 
  AND CB.operador = OP.cod 
  AND CB.tp = '0'
  AND CB.data BETWEEN ? AND ?
  AND CB.empresa = ?
GROUP BY RD.finalizadora, FZ.dsc, OP.nome`
}
