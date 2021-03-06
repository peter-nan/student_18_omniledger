package ch.epfl.dedis;

import ch.epfl.dedis.lib.Roster;
import ch.epfl.dedis.lib.ServerIdentity;
import com.moandjiezana.toml.Toml;

import java.net.URI;
import java.net.URISyntaxException;
import java.util.ArrayList;
import java.util.List;

public class Local {
    public static String rosterStr = "[[servers]]\n" +
            "  Address = \"tcp://localhost:7002\"\n" +
            "  Suite = \"Ed25519\"\n" +
            "  Public = \"a863cf64422ab15f405369134cd057f99e2b40cb45afe7848dde11f34853f708\"\n" +
            "  Description = \"Conode_1\"\n" +
            "[[servers]]\n" +
            "  Address = \"tcp://localhost:7004\"\n" +
            "  Suite = \"Ed25519\"\n" +
            "  Public = \"4706d99de05a58179ccc11ea3c452d9e44b43290de696f83f0fbc8ae26b6679a\"\n" +
            "  Description = \"Conode_2\"\n" +
            "[[servers]]\n" +
            "  Address = \"tcp://localhost:7006\"\n" +
            "  Suite = \"Ed25519\"\n" +
            "  Public = \"4c4d5dd6fa750d5fb32f005b0a357a39d3886454d9fe63255a89ef0542f835d9\"\n" +
            "  Description = \"Conode_3\"\n";
    public static Roster roster = Roster.FromToml(rosterStr);
}
