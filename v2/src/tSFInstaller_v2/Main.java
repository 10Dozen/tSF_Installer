package tSFInstaller_v2;

import javafx.application.Application;
import javafx.geometry.Insets;
import javafx.geometry.Pos;
import javafx.scene.Scene;
import javafx.scene.control.Label;
import javafx.scene.control.TextField;
import javafx.scene.layout.ColumnConstraints;
import javafx.scene.layout.GridPane;
import javafx.scene.text.Font;
import javafx.scene.text.FontWeight;
import javafx.scene.text.Text;
import javafx.stage.Stage;

public class Main extends Application {

    @Override
    public void start(Stage primaryStage) {
        primaryStage.setTitle("tSF Installer (v2.0)");

        GridPane grid = new GridPane();
        grid.setAlignment(Pos.CENTER);
        grid.setHgap(20);
        grid.setVgap(10);
        grid.setPadding(new Insets(25,25,25,25));

        ColumnConstraints column1 = new ColumnConstraints();
        column1.setPercentWidth(30);
        ColumnConstraints column2 = new ColumnConstraints();
        column2.setPercentWidth(70);
        grid.getColumnConstraints().addAll(column1, column2);


        Text title = new Text("tS Framework Installer");
        title.setFont(Font.font("Tahoma", FontWeight.BOLD, 20));


        grid.add(title,0,0);

        Label path = new Label("Installation folder");
        grid.add(path, 0,1);

        TextField pathField = new TextField();
        grid.add(pathField, 1,1);




        Scene scene = new Scene(grid, 640, 480);
        primaryStage.setScene(scene);
        primaryStage.centerOnScreen();
        primaryStage.show();
    }


    public static void main(String[] args) {
        launch(args);
    }
}
