package main;

import java.util.LinkedList;
import java.util.Random;
import java.util.concurrent.CyclicBarrier;
import java.util.concurrent.Semaphore;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

import java.lang.*;

class StringManipulator extends Thread {
    String str = "";
    CyclicBarrier barrier;
    boolean stopped = false;
    boolean equal = false;
    int count = 0;

    StringManipulator(CyclicBarrier b) {
        barrier = b;
        char[] alphabet = {'A', 'B', 'C', 'D'};

        Random r = new Random();
        int size = 10;
        for (int i = 0; i < size; i++) {
            str += alphabet[r.nextInt(4)];
        }
        System.out.println(Thread.currentThread().getName() + " " + str);
    }

    public void run() {
        while (true) {
            for (int i = 0; i < str.length(); i++) {
                if (str.charAt(i) == 'A') {
                    count++;
                }
                if (str.charAt(i) == 'B') {
                    count++;
                }
            }

            try {
                barrier.await();
                barrier.await();
            } catch (Exception err) {
                System.err.println(err);
            }

            if (stopped) {
                System.out.println(str);
                return;
            }

            Random r = new Random();
            int pos = r.nextInt(str.length());
            if (str.charAt(pos) == 'A') {
                str = str.substring(0, pos) + 'C' + str.substring(pos + 1);
            } else 
            if (str.charAt(pos) == 'B') {
                str = str.substring(0, pos) + 'D' + str.substring(pos + 1);
            } else 
            if (str.charAt(pos) == 'C') {
                str = str.substring(0, pos) + 'A' + str.substring(pos + 1);
            } else 
            if (str.charAt(pos) == 'D') {
                str = str.substring(0, pos) + 'B' + str.substring(pos + 1);
            }

            System.out.println(Thread.currentThread().getName() + " " + str);
        }
    }
}

class Resolver extends Thread {
    StringManipulator[] manipulators;
    CyclicBarrier barrier;

    Resolver(CyclicBarrier b, StringManipulator[] m) {
        barrier = b;
        manipulators = m;
    }

    public void run() {
        while (true) {
            try {
                barrier.await();
                
                if ((manipulators[0].count == manipulators[1].count && manipulators[1].count == manipulators[2].count) ||
                (manipulators[1].count == manipulators[2].count && manipulators[2].count == manipulators[3].count) ||
                (manipulators[0].count == manipulators[1].count && manipulators[1].count == manipulators[3].count) ||
                (manipulators[0].count == manipulators[2].count && manipulators[2].count == manipulators[3].count)) {
                    for (int i = 0; i < manipulators.length; i++) {
                        manipulators[i].stopped = true;
                    }

                    barrier.await();
                    return;
                }

                barrier.await();
            } catch (Exception err) {
                System.err.println(err);
            }
        }
    }
}

public class Main {
    public static void main(String[] args) {
        StringManipulator[] m = new StringManipulator[4];
        CyclicBarrier b = new CyclicBarrier(5);

        for (int i = 0; i < 4; i++) {
            m[i] = new StringManipulator(b);
            m[i].start();
        }

        Resolver r = new Resolver(b, m);
        r.start();
        System.out.println("start");
    }
}