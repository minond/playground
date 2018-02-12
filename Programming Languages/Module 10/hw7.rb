# University of Washington, Programming Languages, Homework 7, hw7.rb
# (See also ML code)

# a little language for 2D geometry objects

# each subclass of GeometryExpression, including subclasses of GeometryValue,
#  needs to respond to messages preprocess_prog and eval_prog
#
# each subclass of GeometryValue additionally needs:
#   * shift
#   * intersect, which uses the double-dispatch pattern
#   * intersectPoint, intersectLine, and intersectVerticalLine for
#       for being called by intersect of appropriate clases and doing
#       the correct intersection calculuation
#   * (We would need intersectNoPoints and intersectLineSegment, but these
#      are provided by GeometryValue and should not be overridden.)
#   *  intersectWithSegmentAsLineResult, which is used by
#      intersectLineSegment as described in the assignment
#
# you can define other helper methods, but will not find much need to

# Note: geometry objects should be immutable: assign to fields only during
#       object construction

# Note: For eval_prog, represent environments as arrays of 2-element arrays
# as described in the assignment

# 3. Complete the Ruby implementation except for intersection, which means skip
# for now additions to the Intersect class and, more importantly, methods
# related to intersection in other classes. Do not modify the code given to
# you. Follow this approach:
#
#  - Every subclass of GeometryExpression should have a preprocess_prog method
#    that takes no arguments and returns the geometry object that is the result
#    of preprocessing self. To avoid mutation, return a new instance of the
#    same class unless it is trivial to determine that self is already an
#    appropriate result.
#
#  - Every subclass of GeometryExpression should have an eval_prog method that
#    takes one argu- ment, the environment, which you should represent as an
#    array whose elements are two-element arrays: a Ruby string (the variable
#    name) in index 0 and an object that is a value in our language in index 1.
#    As in any interpreter, pass the appropriate environment when evaluating
#    subexpres- sions. (This is fairly easy since we do not have closures.) To
#    make sure you handle both scope and shadowing correctly:
#
#    - Do not ever mutate an environment; create a new environment as needed
#      instead. Be careful what methods you use on arrays to avoid mutation.
#
#    - The eval_prog method in Var is given to you. Make sure the environments
#      you create work correctly with this definition.
#
#    The result of eval_prog is the result of “evaluating the expression
#    represented by self,” so, as we expect with OOP style, the cases of ML’s
#    eval_prog are spread among our classes, just like with preprocess_prog.
#
#  - Every subclass of GeometryValue should have a shift method that takes two
#    arguments dx and dy and returns the result of shifting self by dx and dy.
#    In other words, all values in the language “know how to shift themselves
#    to create new objects.” Hence the eval_prog method in the Shift class
#    should be very short.
#
#  - Remember you should not use any method like is_a?, instance_of?, class,
#    etc.
#
#  - Analogous to SML, an overall programe would be evaluated via
#    e.preprocess_prog.eval_prog [] (notice we use an array for the
#    environment).

class GeometryExpression
  # do *not* change this class definition
  Epsilon = 0.00001
end

class GeometryValue
  # do *not* change methods in this class definition
  # you can add methods if you wish

  private
  # some helper methods that may be generally useful
  def real_close(r1,r2)
    (r1 - r2).abs < GeometryExpression::Epsilon
  end
  def real_close_point(x1,y1,x2,y2)
    real_close(x1,x2) && real_close(y1,y2)
  end
  # two_points_to_line could return a Line or a VerticalLine
  def two_points_to_line(x1,y1,x2,y2)
    if real_close(x1,x2)
      VerticalLine.new x1
    else
      m = (y2 - y1).to_f / (x2 - x1)
      b = y1 - m * x1
      Line.new(m,b)
    end
  end

  public
  # we put this in this class so all subclasses can inherit it:
  # the intersection of self with a NoPoints is a NoPoints object
  def intersectNoPoints np
    np # could also have NoPoints.new here instead
  end

  # we put this in this class so all subclasses can inhert it:
  # the intersection of self with a LineSegment is computed by
  # first intersecting with the line containing the segment and then
  # calling the result's intersectWithSegmentAsLineResult with the segment
  def intersectLineSegment seg
    line_result = intersect(two_points_to_line(seg.x1,seg.y1,seg.x2,seg.y2))
    line_result.intersectWithSegmentAsLineResult seg
  end
end

class NoPoints < GeometryValue
  # do *not* change this class definition: everything is done for you
  # (although this is the easiest class, it shows what methods every subclass
  # of geometry values needs)
  # However, you *may* move methods from here to a superclass if you wish to

  # Note: no initialize method only because there is nothing it needs to do
  def eval_prog env
    self # all values evaluate to self
  end
  def preprocess_prog
    self # no pre-processing to do here
  end
  def shift(dx,dy)
    self # shifting no-points is no-points
  end
  def intersect other
    other.intersectNoPoints self # will be NoPoints but follow double-dispatch
  end
  def intersectPoint p
    self # intersection with point and no-points is no-points
  end
  def intersectLine line
    self # intersection with line and no-points is no-points
  end
  def intersectVerticalLine vline
    self # intersection with line and no-points is no-points
  end
  # if self is the intersection of (1) some shape s and (2)
  # the line containing seg, then we return the intersection of the
  # shape s and the seg.  seg is an instance of LineSegment
  def intersectWithSegmentAsLineResult seg
    self
  end
end


class Point < GeometryValue
  # *add* methods to this class -- do *not* change given code and do not
  # override any methods

  # Note: You may want a private helper method like the local
  # helper function inbetween in the ML code
  attr_reader :x, :y

  def initialize(x,y)
    @x = x
    @y = y
  end

  def eval_prog env
    self
  end

  def preprocess_prog
    self
  end
end

class Line < GeometryValue
  # *add* methods to this class -- do *not* change given code and do not
  # override any methods
  attr_reader :m, :b

  def initialize(m,b)
    @m = m
    @b = b
  end

  def eval_prog env
    self
  end

  def preprocess_prog
    self
  end
end

class VerticalLine < GeometryValue
  # *add* methods to this class -- do *not* change given code and do not
  # override any methods
  attr_reader :x

  def initialize x
    @x = x
  end

  def eval_prog env
    self
  end

  def preprocess_prog
    self
  end
end

class LineSegment < GeometryValue
  # *add* methods to this class -- do *not* change given code and do not
  # override any methods
  # Note: This is the most difficult class.  In the sample solution,
  #  preprocess_prog is about 15 lines long and
  # intersectWithSegmentAsLineResult is about 40 lines long
  attr_reader :x1, :y1, :x2, :y2
  def initialize (x1,y1,x2,y2)
    @x1 = x1
    @y1 = y1
    @x2 = x2
    @y2 = y2
  end

  def eval_prog env
    self
  end

  # No LineSegment anywhere in the expression has endpoints that are the same
  # as (i.e., real close to) each other. Such a line-segment should be replaced
  # with the appropriate Point. For example in ML syntax,
  # LineSegment(3.2,4.1,3.2,4.1) should be replaced with Point(3.2,4.1).
  #
  # Every LineSegment has its first endpoint (the first two real values in SML)
  # to the left (lower x-value) of the second endpoint. If the x-coordinates of
  # the two endpoints are the same (real close), then the LineSegment has its
  # first endpoint below (lower y-value) the second endpoint. For any
  # LineSegment not meeting this requirement, replace it with a LineSegment
  # with the same endpoints reordered.
  def preprocess_prog
    if real_close(x1, x2) and real_close(y1, y2)
      Point.new(x1, y2)
    elsif y2 < y1
      LineSegment.new(x2, y2, x1, y1)
    else
      self
    end
  end
end

# Note: there is no need for getter methods for the non-value classes

class Intersect < GeometryExpression
  # *add* methods to this class -- do *not* change given code and do not
  # override any methods
  def initialize(e1,e2)
    @e1 = e1
    @e2 = e2
  end

  def preprocess_prog
    Intersect.new(@e1.preprocess_prog, @e2.preprocess_prog)
  end
end

class Let < GeometryExpression
  # *add* methods to this class -- do *not* change given code and do not
  # override any methods
  # Note: Look at Var to guide how you implement Let
  def initialize(s,e1,e2)
    @s = s
    @e1 = e1
    @e2 = e2
  end

  def eval_prog env
    @e2.eval_prog ([@s, @e1.eval_prog(env)] + env)
  end

  def preprocess_prog
    Let.new(@s, @e1.preprocess_prog, @e2.preprocess_prog)
  end
end

class Var < GeometryExpression
  # *add* methods to this class -- do *not* change given code and do not
  # override any methods
  def initialize s
    @s = s
  end

  def eval_prog env # remember: do not change this method
    pr = env.assoc @s
    raise "undefined variable" if pr.nil?
    pr[1]
  end

  def preprocess_prog
    self
  end
end

class Shift < GeometryExpression
  # *add* methods to this class -- do *not* change given code and do not
  # override any methods
  def initialize(dx,dy,e)
    @dx = dx
    @dy = dy
    @e = e
  end

  def eval_prog env
    @e.eval_prog(env).shift(@dx, @dy)
  end

  def preprocess_prog
    Shift.new(@dx, @dy, @e.preprocess_prog)
  end
end
