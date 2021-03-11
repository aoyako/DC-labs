package main;

import java.util.ArrayDeque;
import java.util.concurrent.*;
import java.util.concurrent.atomic.*;
import java.util.concurrent.locks.*;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.sql.Time;
import java.util.Random;

class ArraySummator extends Thread {
    int[] array = new int[100];
    ArraySummator summator;
    static CyclicBarrier barrier1 = new CyclicBarrier(3);
    static CyclicBarrier barrier2 = new CyclicBarrier(3);
    int sum = 0;
    boolean valid = true;

    public ArraySummator(ArraySummator s) {
        summator = s;
        Random r = new Random();
        array[0] = r.nextInt(100);
    }

    public void run() {
        while (true) {
            try {
                valid = false;
                sum = 0;
                for (int i = 0; i < 100; i++) {
                    sum += array[i];
                }
                System.out.println(Thread.currentThread().getName() + " " + sum);
        
                barrier1.await();
        
                if (summator.sum + 1 < sum) {
                    array[0]--;
                } else
                if (summator.sum > sum) {
                    array[0]++;
                } else
                if (summator.sum == sum) {
                    valid = true;
                }

                barrier2.await();

                if (summator.valid && valid) {
                    System.out.println("stopped");
                    return;
                }

            } catch (Exception err) {
                System.err.println(err);
            }
        }
    }
}

public class Main {
    public static void main(String[] args) {
        ArraySummator s1 = new ArraySummator(null);
        ArraySummator s2 = new ArraySummator(s1);
        ArraySummator s3 = new ArraySummator(s2);
        s1.summator = s3;

        s1.start();
        s2.start();
        s3.start();
    }
}