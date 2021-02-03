package main;

import javax.swing.*;  
import javax.swing.event.*;
import java.awt.event.*;

class GoTo20 implements Runnable {
    private JSlider sl;

    public GoTo20(JSlider sl) {
        this.sl = sl;
    }

    @Override
    public void run() {
        while (true) {
            synchronized(sl) {
                int currentState = sl.getValue();
                if (currentState > 20) {
                    currentState--;
                } else if (currentState < 20) {
                    currentState++;
                }
                sl.setValue(currentState);
            }
            // System.out.println("Hello");

            try {
                Thread.sleep(200/Thread.currentThread().getPriority());
            } catch (Exception err) {

            }
        }
    }
}

class GoTo90 implements Runnable {
    private JSlider sl;

    public GoTo90(JSlider sl) {
        this.sl = sl;
    }

    @Override
    public void run() {
        while (true) {
            synchronized(sl) {
                int currentState = sl.getValue();
                if (currentState > 90) {
                    currentState--;
                } else if (currentState < 90) {
                    currentState++;
                }
                sl.setValue(currentState);
            }
            // System.out.println("Hello");

            try {
                Thread.sleep(200/Thread.currentThread().getPriority());
            } catch (Exception err) {

            }
        }
    }
}

public class Main {
    private static Thread thread1;
    private static Thread thread2;
    public static void main(String[] args) {
        JFrame frame = new JFrame();
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
                
        JSlider positionSlider = new JSlider(0, 100, 50);
        JSlider thread20Priority = new JSlider(1, 10, 5);
        JSlider thread90Priority = new JSlider(1, 10, 5);
        JButton startButton = new JButton("START");
        JLabel thread20Label = new JLabel("Thread1 priority"); 
        JLabel thread90Label = new JLabel("Thread2 priority"); 

        startButton.setSize(80, 80);
        positionSlider.setBounds(130, 100, 300, 40);
        positionSlider.setMajorTickSpacing(10);
        positionSlider.setPaintTicks(true);
        positionSlider.setPaintLabels(true);

        thread20Priority.setBounds(100, 30, 150, 30);
        thread90Priority.setBounds(280, 30, 150, 30);
        thread20Label.setBounds(100, 0, 150, 30);
        thread90Label.setBounds(280, 0, 150, 30);
        thread20Priority.setMajorTickSpacing(1);
        thread20Priority.setPaintTicks(true);
        thread90Priority.setMajorTickSpacing(1);
        thread90Priority.setPaintTicks(true);

        thread20Priority.addChangeListener(
            new ChangeListener() {
                @Override
                public void stateChanged(ChangeEvent evt) {
                    JSlider temp = (JSlider) evt.getSource();
                        thread1.setPriority(temp.getValue());
                 }
            }
        );

        thread90Priority.addChangeListener(
            new ChangeListener() {
                @Override
                public void stateChanged(ChangeEvent evt) {
                    JSlider temp = (JSlider) evt.getSource();
                        thread2.setPriority(temp.getValue());
                 }
            }
        );

        startButton.addActionListener(
            new ActionListener() {
                @Override
                public void actionPerformed(ActionEvent evt) {
                    thread1.start();
                    thread2.start();
                }
            }
        );
                    
        frame.add(positionSlider);
        frame.add(startButton);
        frame.add(thread20Priority);
        frame.add(thread90Priority);
        frame.add(thread20Label);
        frame.add(thread90Label);
                    
        frame.setSize(1000, 500);
        frame.setLayout(null);
        frame.setVisible(true);

        thread1 = new Thread(new GoTo20(positionSlider));
        thread2 = new Thread(new GoTo90(positionSlider));
   }
}