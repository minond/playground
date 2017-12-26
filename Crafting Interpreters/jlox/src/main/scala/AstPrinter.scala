package com.craftinginterpreters.lox

class AstPrinter extends Expr.Visitor[String] {
  def print(expr: Expr): String = {
    expr.accept(this)
  }

  override def visitBinaryExpr(expr: Expr.Binary): String = {
    parenthesize(expr.operator.lexeme, expr.left, expr.right)
  }

  override def visitGroupingExpr(expr: Expr.Grouping): String = {
    parenthesize("group", expr.expression)
  }

  override def visitLiteralExpr(expr: Expr.Literal): String = {
    if (expr.value == null)
      "nil"
    else
      expr.value.toString
  }

  override def visitUnaryExpr(expr: Expr.Unary): String = {
    parenthesize(expr.operator.lexeme, expr.right)
  }

  private def parenthesize(name: String, exprs: Expr*): String = {
    val builder = new StringBuilder()

    builder
      .append("(")
      .append(name)

    exprs.foreach { expr =>
      builder
        .append(" ")
        .append(expr.accept(this))
    }

    builder.append(")")
    return builder.toString
  }
}
