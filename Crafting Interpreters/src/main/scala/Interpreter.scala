package com.craftinginterpreters.lox

import scala.util.{Try, Success, Failure}
import com.craftinginterpreters.lox.TokenType._

class Interpreter extends Expr.Visitor[Any] {
  def interpret(expr: Expr): Unit  = {
    Try { stringify(evaluate(expr)) } match {
      case Success(value) => println(value)
      case Failure(error: RuntimeError) => Main.runtimeError(error)
      case Failure(error) => throw error
    }
  }

  override def visitBinaryExpr(expr: Expr.Binary): Any = {
    val left = evaluate(expr.left)
    val right = evaluate(expr.right)

    expr.operator.ttype match {
      // Arithmetic operator
      case MINUS => asOpNumber(expr.operator, left) - asOpNumber(expr.operator, right)
      case SLASH => asOpNumber(expr.operator, left) / asOpNumber(expr.operator, right)
      case STAR => asOpNumber(expr.operator, left) * asOpNumber(expr.operator, right)
      case PLUS =>
        if (left.isInstanceOf[Double] && right.isInstanceOf[Double])
          left.asInstanceOf[Double] + right.asInstanceOf[Double]
        else if (left.isInstanceOf[String] && right.isInstanceOf[String])
          left.asInstanceOf[String] + right.asInstanceOf[String]
        else
          throw new RuntimeError(expr.operator, "Operands must be two numbers or two strings.")

      // Relational operators
      case GREATER => asOpNumber(expr.operator, left) > asOpNumber(expr.operator, right)
      case GREATER_EQUAL => asOpNumber(expr.operator, left) >= asOpNumber(expr.operator, right)
      case LESS => asOpNumber(expr.operator, left) < asOpNumber(expr.operator, right)
      case LESS_EQUAL => asOpNumber(expr.operator, left) <= asOpNumber(expr.operator, right)

      // Equality operators
      case BANG_EQUAL => !isEqual(left, right)
      case EQUAL_EQUAL => isEqual(left, right)

      case _ => throw new RuntimeError(expr.operator, "Unknown operator.")
    }
  }

  override def visitGroupingExpr(expr: Expr.Grouping): Any = {
    evaluate(expr.expression)
  }

  override def visitLiteralExpr(expr: Expr.Literal): Any = {
    expr.value match {
      case Some(value) => value
      case None => null
    }
  }

  override def visitUnaryExpr(expr: Expr.Unary): Any = {
    val right = evaluate(expr.right)

    expr.operator.ttype match {
      case MINUS => -asOpNumber(expr.operator, right)
      case BANG => !isTruthy(right)
      case _ => null
    }
  }

  private def evaluate(expr: Expr): Any = {
    expr.accept(this)
  }

  private def isTruthy(value: Any): Boolean = {
    if (value == null)
      false
    else if (value.isInstanceOf[Boolean])
      value.asInstanceOf[Boolean]
    else
      true
  }

  private def isEqual(left: Any, right: Any): Boolean = {
    if (left == null && right == null)
      true
    else if (left == null)
      false
    else
      left == right
  }

  private def asOpNumber(op: Token, num: Any): Double = {
    if (num.isInstanceOf[Double])
      num.asInstanceOf[Double]
    else
      throw new RuntimeError(op, "Operand must be a number.")
  }

  private def stringify(value: Any): String = {
    if (value == null) {
      "nil"
    } else if (value.isInstanceOf[Double]) {
      val txt = value.toString

      if (txt.endsWith(".0"))
        txt.stripSuffix(".0")
      else
        txt
    } else {
      value.toString
    }
  }
}
