#  快速将C 结构体转换成 asciiTable 的工具



例子:

```c
typedef struct clusterLink {
    struct B {
        int **e /*sadfsa*/, d /*fasdfas*/, f;/*fff*/
    } *val;
    /* hello world*/struct c{
        int k; /*1*/
        int z1, d1; /*2*/
        double do1; /**/
    }  p;

} clusterLink;
/*

 are you ok ?

*/
```

转换成：

```html
/*
 *        +----------+---------+--------------+
 *        | TYPENAME | VARNAME |    EXTRA     |
 *        +----------+---------+--------------+
 *        |    B*    |   val   |  hello world |
 *        +----------+---------+--------------+
 *        |  int**   |    e    |    sadfsa    |
 *        +----------+---------+--------------+
 *        |   int    |    d    |   fasdfas    |
 *        +          +---------+--------------+
 *        |          |    f    |     fff      |
 *        +----------+---------+--------------+
 */


/*
 *        +----------+---------+-------+
 *        | TYPENAME | VARNAME | EXTRA |
 *        +----------+---------+-------+
 *        |    c     |    p    |       |
 *        +----------+---------+-------+
 *        |   int    |    k    |   1   |
 *        +          +---------+-------+
 *        |          |   z1    |       |
 *        +          +---------+-------+
 *        |          |   d1    |   2   |
 *        +----------+---------+-------+
 *        |  double  |   do1   |       |
 *        +----------+---------+-------+
 */


/*
 *        +-------------+-------------+---------------+
 *        |  TYPENAME   |   VARNAME   |     EXTRA     |
 *        +-------------+-------------+---------------+
 *        | clusterLink | clusterLink |    are you ok |
 *        |             |             |      ?        |
 *        +-------------+-------------+---------------+
 *        |     B*      |     val     |  hello world  |
 *        +-------------+-------------+---------------+
 *        |      c      |      p      |               |
 *        +-------------+-------------+---------------+
 */
```

暂时还不支持解析括号。