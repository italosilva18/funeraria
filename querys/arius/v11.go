package arius

type Arius11 struct{}

func (g *Arius11) TOTAL_VENDIDO_DIA() string {
	return "SELECT SUM(totalValor) AS Valor FROM(SELECT SUM(it.Valor)+SUM(it.JuroCheque)+SUM(it.JuroCartao) AS totalValor FROM retag.cupom cp LEFT JOIN retag.itens it ON it.nroloja = cp.nroloja AND cp.pdv = it.pdv AND cp.DataProc = it.DataProc AND it.NroCupom = cp.NroCupom LEFT JOIN retag.mercador mcd ON it.nroloja = mcd.nroloja AND it.Codigo = mcd.codigoean WHERE cp.nroloja = ? AND cp.tipooperacao = 1 AND cp.FlagFimCupom = 1 AND cp.FlagEstorno = 0 AND cp.dataproc between ? AND ? AND cp.HoraMinSeg BETWEEN ? AND ? AND it.Estornado = 0 GROUP BY cp.pdv, cp.NroCupom) t;"
}

// TICKET_MEDIO_DIA - Alterada 01-02-2019
func (g *Arius11) TICKET_MEDIO_DIA() string {
	return "SELECT Round(SUM(totalValor)/COUNT(nroitens), 2) As MediaValor FROM (SELECT SUM(it.Valor)+SUM(it.JuroCheque)+SUM(it.JuroCartao) AS totalValor, COUNT( * ) as nroItens FROM retag.cupom cp LEFT JOIN retag.itens it ON it.nroloja = cp.nroloja AND cp.pdv = it.pdv AND cp.DataProc = it.DataProc AND it.NroCupom = cp.NroCupom LEFT JOIN retag.mercador mcd ON it.nroloja = mcd.nroloja AND it.Codigo = mcd.codigoean WHERE cp.nroloja = ? AND cp.tipooperacao = 1 AND cp.FlagFimCupom = 1 AND cp.FlagEstorno = 0 AND cp.dataproc BETWEEN ? AND ? AND cp.HoraMinSeg BETWEEN ? AND ? and it.Estornado = 0 GROUP BY cp.pdv, cp.NroCupom) t;"
}

// PRODUTOS_POR_ATENDIMENTO_DIA - Alterada 01-02-2019
func (g *Arius11) PRODUTOS_POR_ATENDIMENTO_DIA() string {
	return "SELECT Round(sum(nroitens)/count(nroitens) , 2) As MediaItens FROM (SELECT SUM(it.Valor)+SUM(it.JuroCheque)+SUM(it.JuroCartao) AS totalValor, COUNT(*) as nroItens FROM retag.cupom cp LEFT JOIN retag.itens it ON it.nroloja = cp.nroloja AND cp.pdv = it.pdv AND cp.DataProc = it.DataProc AND it.NroCupom = cp.NroCupom LEFT JOIN retag.mercador mcd ON it.nroloja = mcd.nroloja AND it.Codigo = mcd.codigoean WHERE cp.nroloja = ? AND cp.tipooperacao = 1 AND cp.FlagFimCupom = 1 AND cp.FlagEstorno = 0 AND cp.dataproc BETWEEN ? AND ? and cp.HoraMinSeg BETWEEN ? AND ? and it.Estornado = 0 GROUP BY cp.pdv, cp.NroCupom) t;"
}

// TOTAL_CUPONS_VALIDOS_DIA - Alterada 01-02-2019
func (g *Arius11) TOTAL_CUPONS_VALIDOS_DIA() string {
	return "SELECT CAST(COUNT(*) AS signed) AS quantidade FROM retag.cupom AS c WHERE c.dataproc BETWEEN ? AND ? AND c.flagestorno = 0 AND c.tipoOperacao = 1 AND c.recebimento = 0 AND c.nroloja = ? AND c.flagfimcupom = 1"
}

// TOTAL_CUPONS_CANCELADO_DIAS - Alterada 01-02-2019
func (g *Arius11) TOTAL_CUPONS_CANCELADOS_DIA() string {
	return "SELECT CAST(COUNT(*) AS signed) AS 'Cupons Cancelados no Mes' FROM retag.cupom AS c WHERE c.dataproc BETWEEN ? AND ? AND c.flagestorno <> 0 AND c.tipoOperacao = 1 AND c.recebimento = 0 AND c.nroloja = ?"
}

// VENDAS_POR_FAIXA_HORARIA - Alterada 01-02-2019
func (g *Arius11) VENDAS_POR_FAIXA_HORARIA() string {
	return "SELECT Hora, sum(totalValor) as Valor, count(nroitens) as NumCupons, sum(nroitens) as NumItens, Round(sum(nroitens)/count(nroitens) , 2) As MediaItens, Round(sum(totalValor)/ count(nroitens), 2) As MediaValor FROM (SELECT SUM(it.Valor)+SUM(it.JuroCheque)+SUM(it.JuroCartao) AS totalValor, substr(cp.HoraMinSeg,12,2) as hora, COUNT(*) as nroItens FROM retag.cupom cp LEFT JOIN retag.itens it ON it.nroloja = cp.nroloja AND cp.pdv = it.pdv AND cp.DataProc = it.DataProc AND it.NroCupom = cp.NroCupom LEFT JOIN retag.mercador mcd ON it.nroloja = mcd.nroloja AND it.Codigo = mcd.codigoean WHERE cp.nroloja = ? AND cp.tipooperacao = 1 AND cp.FlagFimCupom = 1 AND cp.FlagEstorno = 0 AND cp.dataproc = ? and cp.HoraMinSeg BETWEEN ? AND ? and it.Estornado = 0 GROUP BY cp.pdv, cp.NroCupom) t GROUP BY Hora;"
}

// PRODUTOS_MAIS_VENDIDOS_DIA - Alterada 01-02-2019
func (g *Arius11) PRODUTOS_MAIS_VENDIDO_DIA() string {
	return "SELECT a.codigo, b.descricao, SUM(a.quantidade) AS Quantidade, SUM(a.valor) as Valor FROM retag.itens AS a LEFT JOIN retag.cupom AS e ON a.nroloja = e.nroloja AND a.nrocupom = e.nrocupom AND a.dataproc = e.dataproc AND a.pdv = e.pdv LEFT JOIN retag.mercador AS b ON a.nroloja = b.nroloja AND a.codigo = b.codigoean WHERE a.nroloja = ? AND e.dataproc = ? AND a.estornado = 0 AND e.flagfimcupom = 1 AND e.flagestorno = 0 AND e.recebimento = 0 GROUP BY a.codigo, b.descricao ORDER BY valor DESC LIMIT 0, 200"
}

// PRODUTOS_ESTATISTICA - Alterada 01-02-2019
func (g *Arius11) PRODUTOS_ESTATISTICA() string {
	return "SELECT a.codigo, b.descricao, a.ValorUnitario, SUM(a.Quantidade) AS Quantidade, SUM(a.valor) AS valor FROM retag.itens AS a LEFT JOIN retag.cupom AS e ON a.nroloja = e.nroloja AND a.nrocupom = e.nrocupom AND a.dataproc = e.dataproc LEFT JOIN retag.mercador AS b ON a.nroloja = b.nroloja AND a.codigo = b.codigoean WHERE a.nroloja = ? AND a.dataproc = date(?) AND a.estornado = 0 AND e.flagfimcupom = 1 AND e.flagestorno = 0 AND e.recebimento = 0 GROUP BY a.codigo, b.descricao ORDER BY valor DESC LIMIT 8000"
}

// TOTAL_VENDIDO_POR_MODALIDADES - Alterada 01-02-2019
func (g *Arius11) TOTAL_VENDIDO_POR_MODALIDADES() string {
	return "SELECT f.CodFpagto AS meio_pagto, a.descricao, count(1) AS qtd, SUM(f.valor-f.troco/100) AS valor FROM retag.fpagtoCupom AS f LEFT JOIN retag.cupom AS c ON f.nroloja = c.nroloja AND f.dataproc = c.dataproc AND f.pdv = c.pdv AND f.nrocupom = c.nrocupom LEFT JOIN controle.meiospagto AS a ON f.nroloja = a.nroloja and f.codfpagto = a.codigo WHERE f.nroloja = ? AND f.dataproc = ? AND c.flagfimcupom = 1 and f.flgtroca = 0 AND (c.tipooperacao = 1 OR c.tipooperacao = 6) GROUP BY a.codigo"
}

// ESTORNO_DE_CUPONS_POR_OPERADOR_DIA - Alterada 01-02-2019
func (g *Arius11) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA() string {
	return "SELECT o.pdv, o.DataProc AS dataproc, o.hora, LEFT( u.nome, 35 ) AS operador, LEFT( us.nome, 35 ) AS supervisor, o.valor, o.nrocupom FROM retag.ocorrencias o LEFT JOIN controle.usuarios u ON o.nroloja = u.nroloja AND o.operador = u.codigo LEFT JOIN controle.usuarios us ON o.nroloja = us.nroloja AND o.supervisor = us.codigo LEFT JOIN retag.mercador m ON m.nroloja = o.nroloja AND o.codigoean = m.codigoean WHERE o.nroloja = ? AND dataproc BETWEEN ? AND ? AND o.descricao = 'ESTORNO CUPOM'"
}

// DESCONTOS_DE_ITENS_DIA - Alterada 01-02-2019
func (g *Arius11) DESCONTOS_DE_ITENS_DIA() string {
	return "SELECT o.DataProc AS dataproc, o.hora, LEFT( u.nome, 35 ) AS operador, LEFT( us.nome, 35 ) AS supervisor, o.valor, o.nrocupom, o.item, o.codigoean, m.descricao AS desc_item FROM retag.ocorrencias o LEFT JOIN controle.usuarios u ON o.nroloja = u.nroloja AND o.operador = u.codigo LEFT JOIN controle.usuarios us ON o.nroloja = us.nroloja AND o.supervisor = us.codigo LEFT JOIN retag.mercador m ON m.nroloja = o.nroloja AND o.codigoean = m.codigoean WHERE o.nroloja = ? AND dataproc BETWEEN ? AND ? AND o.descricao = 'DESCONTO ITEM' and o.supervisor <> '000000' ORDER BY o.DataProc, o.descricao,  o.operador"
}

// VENDAS_POR_SECOES_DIA - Alterada 01-02-2019
func (g *Arius11) VENDAS_POR_SECOES_DIA() string {
	return "SELECT i.depto AS Cod_secao, s.descricao, SUM(i.quantidade) AS qtd, SUM(i.valor+i.jurocheque+i.jurocartao) AS valor FROM retag.itens AS i LEFT JOIN retag.cupom AS c ON i.nroloja = c.nroloja and i.nrocupom = c.nrocupom AND i.pdv = c.pdv And i.dataproc = c.dataproc LEFT JOIN controle.secoes AS s ON i.depto = s.codigo WHERE i.nroloja = ? AND i.dataproc = ? AND i.estornado = 0 AND c.flagfimcupom = 1 AND c.flagestorno = 0 AND c.recebimento = 0 GROUP BY Cod_secao"
}

// TOTAL_CUPONS_POR_OPERADOR_DIA - Alterada 01-02-2019
func (g *Arius11) TOTAL_CUPONS_POR_OPERADOR_DIA() string {
	return "SELECT CAST(COUNT(*) AS signed) as Qtd_Cupom, LEFT( u.nome, 35 ) AS operador FROM retag.cupom AS o LEFT JOIN controle.usuarios AS u ON o.nroloja = u.nroloja WHERE o.operador = u.codigo AND dataproc BETWEEN ? AND ? AND o.tipooperacao = 1 AND o.FlagFimCupom = 1 AND o.FlagEstorno = 0 AND o.recebimento = 0 AND o.nroloja = ? GROUP BY operador ORDER BY Qtd_cupom DESC;"
}

// VENDAS_POR_VENDEDOR_DIA - Alterada 01-02-2019
func (g *Arius11) VENDAS_POR_VENDEDOR_DIA() string {
	return "SELECT Data, case  WHEN vendedor is null then 'Sem Vendedor' else `vendedor` END AS vendedor, sum(totalValor) as Valor, count(nroitens) as NumCupons, sum(nroitens) as NumItens, Round(sum(nroitens)/count(nroitens) , 2) As MediaItens, Round(sum(totalValor)/ count(nroitens), 2) As MediaValor FROM (SELECT SUM(it.Valor)+SUM(it.JuroCheque)+SUM(it.JuroCartao) AS totalValor, cp.Dataproc as Data, u.nome as Vendedor, COUNT(*) as nroItens FROM retag.cupom cp LEFT JOIN retag.itens it ON it.nroloja = cp.nroloja AND cp.pdv = it.pdv AND cp.DataProc = it.DataProc AND it.NroCupom = cp.NroCupom LEFT JOIN controle.usuarios u ON it.nroloja = u.nroloja AND it.vendedor = u.codigo LEFT JOIN retag.mercador mcd ON it.nroloja = mcd.nroloja AND  it.Codigo = mcd.codigoean WHERE cp.nroloja = ? AND cp.tipooperacao = 1 AND cp.FlagFimCupom = 1 AND cp.FlagEstorno = 0 AND cp.dataproc BETWEEN ? AND ? and cp.HoraMinSeg BETWEEN ? AND ? and it.Estornado = 0 GROUP BY vendedor, cp.pdv, cp.NroCupom) t group BY vendedor order by valor DESC limit 0, 10"
}

// MODALIDADES_PAGAMENTO_OPERADOR_DIA - Alterada 01-02-2019
func (g *Arius11) MODALIDADES_PAGAMENTO_OPERADOR_DIA() string {
	return "SELECT LEFT( u.nome, 35 ) AS operador, f.CodFpagto AS meio_pagto, a.descricao, COUNT(1) AS qtd, SUM(f.valor-f.troco/100) AS valor FROM retag.fpagtoCupom AS f LEFT JOIN retag.cupom AS c ON f.nroloja = c.nroloja AND f.dataproc = c.dataproc AND f.pdv = c.pdv AND f.nrocupom = c.nrocupom LEFT JOIN controle.meiospagto AS a ON f.nroloja = a.nroloja and f.codfpagto = a.codigo LEFT JOIN controle.usuarios AS u ON c.nroloja = u.nroloja WHERE c.operador = u.codigo AND f.nroloja = ? AND f.dataproc BETWEEN ? AND ? AND c.flagfimcupom = 1 and f.flgtroca = 0 AND (c.tipooperacao = 1 OR c.tipooperacao = 6) GROUP BY a.codigo, operador ORDER BY operador;"
}
