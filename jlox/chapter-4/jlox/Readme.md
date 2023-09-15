# CraftingInterpreters

## Useful commands

Create chapter

```bash
mvn archetype:generate -DgroupId=com.craftinginterpreters.lox -DartifactId=jlox -DarchetypeArtifactId=maven-archetype-quickstart -DarchetypeVersion=1.4 -DinteractiveMode=false
```

Build package

```bash
mvn package
```

Run the code

```bash
java -cp target/jlox-1.0-SNAPSHOT.jar com.craftinginterpreters.lox.Lox
```

or

```bash
java -cp target/jlox-1.0-SNAPSHOT.jar com.craftinginterpreters.tools.GenerateAst
```
