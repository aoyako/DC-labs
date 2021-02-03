package main;

import javax.swing.*;  
import javax.swing.event.*;
import java.awt.event.*;
import java.util.concurrent.*;

class GoToPos extends Thread {
    private JSlider sl;
    private boolean dead = false;
    private int pos;

    public GoToPos(JSlider sl, int pos) {
        this.sl = sl;
        this.pos = pos;
    }

    public void kill() {
        dead = true;
    }

    @Override
    public void run() {
        while (true) {
            if (dead) {
                return;
            }

            synchronized(sl) {
                int currentState = sl.getValue();
                if (currentState > pos) {
                    currentState--;
                } else if (currentState < pos) {
                    currentState++;
                }
                sl.setValue(currentState);
            }

            try {
                Thread.sleep(200/Thread.currentThread().getPriority());
            } catch (Exception err) {

            }
        }
    }
}

public class Main {
    private static ThreadWrapper threadWrapper1 = new ThreadWrapper(1);
    private static ThreadWrapper threadWrapper2 = new ThreadWrapper(2);

    private static Semaphore sem = new Semaphore(1);
    private static int workingThreadID = 0;

    public static void main(String[] args) {
        JFrame frame = new JFrame();
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
                
        JSlider positionSlider = new JSlider(0, 100, 50);
        JButton start1 = new JButton("START1");
        JButton start2 = new JButton("START2");
        JButton stop1 = new JButton("STOP1");
        JButton stop2 = new JButton("STOP2");

        positionSlider.setBounds(350, 100, 300, 40);
        positionSlider.setMajorTickSpacing(10);
        positionSlider.setPaintTicks(true);
        positionSlider.setPaintLabels(true);

        start1.setBounds(0, 0, 100, 80);
        stop1.setBounds(0, 100, 100, 80);
        start2.setBounds(120, 0, 100, 80);
        stop2.setBounds(120, 100, 100, 80);

        start1.addActionListener(generateStartButtonListener(threadWrapper1, positionSlider, 20, stop2, frame));
        start2.addActionListener(generateStartButtonListener(threadWrapper2, positionSlider, 90, stop1, frame));

        stop1.addActionListener(generateStopButtonListener(threadWrapper1, stop2));
        stop2.addActionListener(generateStopButtonListener(threadWrapper2, stop1));

        frame.add(positionSlider);
        frame.add(start1);
        frame.add(start2);
        frame.add(stop1);
        frame.add(stop2);

                    
        frame.setSize(1000, 500);
        frame.setLayout(null);
        frame.setVisible(true);

    }
   
    protected static ActionListener generateStartButtonListener(ThreadWrapper t, JSlider s, int postition, JButton blockedStop, JFrame frame) {
        return new ActionListener() {
            private ThreadWrapper tw = t;

            @Override
            public void actionPerformed(ActionEvent evt) {
                try {
                    boolean res = sem.tryAcquire();
                    if (!res) {
                        JOptionPane.showMessageDialog(frame, "In use by thread");
                        return;
                    }
                } catch (Exception err) {
                    System.out.println(err);
                }

                workingThreadID = t.id;
                blockedStop.setEnabled(false);

                tw.thread = new GoToPos(s, postition);
                tw.thread.start();

                if (t.id == 1) {
                    tw.thread.setPriority(1);
                } else {
                    tw.thread.setPriority(10);
                }
            }
        };
    }

    protected static ActionListener generateStopButtonListener(ThreadWrapper t, JButton blockedStop) {
        return new ActionListener() {
            private ThreadWrapper tw = t;

            @Override
            public void actionPerformed(ActionEvent evt) {
                if (workingThreadID != tw.id) {
                    return;
                }

                blockedStop.setEnabled(true);
                tw.thread.kill();

                workingThreadID = 0;

                try {
                    sem.release();
                } catch (Exception err) {
                    System.out.println(err);
                }
            }
        };
    }
}

class ThreadWrapper {
    public ThreadWrapper(int id) {
        this.id = id;
    }

    public GoToPos thread;
    public int id;
}