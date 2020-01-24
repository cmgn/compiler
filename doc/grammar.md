    program
      | {statement}

    statement
      | "{" {statement} "}"
      | "if" expression statement ["else" statement]
      | "while" expression statement
      | "var" identifier type ";"
      | expression "=" expression ";"
      | expression ";"
      | ";"

    type
      | "int"
      | "array" "(" integer ")" "of" type
      | "ptr" "to" type
      | identifier

    expression
      | equality

    equality
      | equality "=" comparison
      | comparison

    comparison
      | summation "<" summation
      | summation ">" summation
      | summation

    summation
      | summation "+" product
      | summation "-" product
      | product

    product
      | product "*" subscript
      | product "/" subscript
      | subscript

    subscript
      | subscript "[" expression "]"
      | terminal

    terminal
      | integer
      | identifier
      | "(" expression ")"
      | "&" terminal
      | "*" terminal
      | "-" terminal
