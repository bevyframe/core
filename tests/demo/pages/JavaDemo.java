import me.islekcaganmert.bevyframe.*;

public class JavaDemo {
    public static void main(String[] args) {
        Context context = new Context();
        if (context.getMethod().equals("GET"))
            get(context).respond();
    }

    public static Response get(Context context) {
        return new Response("Hello, World!".getBytes());
    }
}