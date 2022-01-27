import com.signadot.ApiClient;
import com.signadot.ApiException;
import com.signadot.api.WorkspacesApi;
import com.signadot.model.*;
import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import org.apache.commons.lang3.RandomStringUtils;
import org.testng.annotations.AfterSuite;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.Test;

import java.util.List;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;

public class RouteServiceTest {

    public static final String ORG_NAME = "signadot";
    public static final String HOTROD = "hotrod";
    public static final String SIGNADOT_API_KEY = System.getenv("SIGNADOT_API_KEY");
    public static final String ROUTE_SERVICE_IMAGE_NAME = System.getenv("ROUTE_SERVICE_IMAGE_NAME");

    ApiClient apiClient;
    WorkspacesApi workspacesApi;
    CreateWorkspaceResponse response;
    String workspaceID;

    @BeforeSuite
    public void createWorkspace() throws ApiException, InterruptedException {
        apiClient = new ApiClient().addDefaultHeader("Accept-Version", "20211220");
        apiClient.setApiKey(SIGNADOT_API_KEY);
        workspacesApi = new WorkspacesApi(apiClient);

        String workspaceName = String.format("test-ws-%s", RandomStringUtils.randomAlphanumeric(5));
        WorkspaceFork routeFork = new WorkspaceFork()
                .forkOf(new ForkOf().kind("Deployment").namespace(HOTROD).name("route"))
                .customizations(new WorkspaceCustomizations()
                        .addEnvItem(new EnvOp().name("abc").value("def").operation("upsert"))
                        .addImagesItem(new Image().image(ROUTE_SERVICE_IMAGE_NAME)))
                .addEndpointsItem(new ForkEndpoint().name("route").port(8083).protocol("http"));

        CreateWorkspaceRequest request = new CreateWorkspaceRequest()
                .cluster("signadot-staging")
                .name(workspaceName)
                .description("test workspace created using signadot-sdk")
                .addForksItem(routeFork);

        response = workspacesApi.createNewWorkspace(ORG_NAME, request);

        workspaceID = response.getWorkspaceID();
        if (workspaceID == null || workspaceID == "") {
            throw new RuntimeException("Workspace ID not set");
        }

        List<PreviewEndpoint> endpoints = response.getPreviewEndpoints();
        if (endpoints.size() == 0) {
            throw new RuntimeException("preview URL not generated");
        }

        PreviewEndpoint endpoint = null;
        for (PreviewEndpoint ep: endpoints) {
            if ("route".equals(ep.getName())) {
                endpoint = ep;
                break;
            }
        }
        if (endpoint == null) {
            throw new RuntimeException("No matching endpoint found");
        }

        // set the base URL for tests
        RestAssured.baseURI = endpoint.getPreviewURL();

        // Check for workspace readiness
        while (!workspacesApi.getWorkspaceReady(ORG_NAME, workspaceID).isReady()) {
            Thread.sleep(5000);
        };
    }

    @Test
    public void testETANotNegative() {
        given().
                header("signadot-api-key", SIGNADOT_API_KEY).
                when().
                get("/route?pickup=123&dropoff=456").
                then().
                statusCode(200).
                assertThat().body("ETA", greaterThan(Long.valueOf(-1)));
    }

    @Test
    public void testPickupDropOffHasValue() {
        given().
                header("signadot-api-key", SIGNADOT_API_KEY).
                when().
                get("/route?pickup=123&dropoff=456").
                then().
                statusCode(200).
                assertThat().
                body("Pickup", not(emptyOrNullString())).
                body("Dropoff", not(emptyOrNullString()));
    }

    @Test
    public void testStatusCode200() {
        given().
                header("signadot-api-key", SIGNADOT_API_KEY).
                when().
                get("/route?pickup=123&dropoff=456").
                then().
                statusCode(200);
    }

    @Test
    public void testNoQueryParams() {
        given().
                header("signadot-api-key", SIGNADOT_API_KEY).
                when().
                get("/route").
                then().
                statusCode(400).
                contentType(ContentType.TEXT).
                body(containsString("Missing required 'pickup' parameter"));
    }

    @Test
    public void testRequirePickupQueryParam() {
        given().
                header("signadot-api-key", SIGNADOT_API_KEY).
                when().
                get("/route?dropoff=456").
                then().
                statusCode(400).
                contentType(ContentType.TEXT).
                body(containsString("Missing required 'pickup' parameter"));
    }

    @Test
    public void testRequireDropoffQueryParam() {
        given().
                header("signadot-api-key", SIGNADOT_API_KEY).
                when().
                get("/route?pickup=577,322").
                then().
                statusCode(400).
                contentType(ContentType.TEXT).
                body(containsString("Missing required 'dropoff' parameter"));
    }

    @AfterSuite
    public void deleteWorkspace() throws ApiException {
        workspacesApi.deleteWorkspaceById(ORG_NAME, workspaceID);
    }
}
